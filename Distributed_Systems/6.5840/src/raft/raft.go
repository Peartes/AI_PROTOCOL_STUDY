package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	//	"bytes"

	"bytes"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"6.5840/labgob"
	"6.5840/labrpc"
)

// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 3D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 3D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()
	timeout   time.Time
	state     State

	// Your data here (3A, 3B, 3C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	// persistent state
	currentTerm int // current term
	votedFor    int
	logs        []LogEntries

	leaderId int // the index of the leader in peers
}

type State string

var Leader State = "Leader"
var Follower State = "Follower"
var Candidate State = "Candidate"

type LogEntries struct {
	Command string
	Term    int
}

type AppendEntry struct {
	Term         int          // leader’s term
	LeaderId     int          // so follower can redirect clients
	PrevLogIndex int          // index of log entry immediately preceding new ones
	PrevLogTerm  int          // term of prevLogIndex entry
	Entries      []LogEntries // log entries to store (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int          // leader’s commitIndex
}

type AppendEntryReply struct {
	Term    int  // for leader to update himself
	Success bool // true if follower contained entry matching prevLogIndex and prevLogTerm
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	// Your code here (3A).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.currentTerm, rf.state == Leader
}

// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
// before you've implemented snapshots, you should pass nil as the
// second argument to persister.Save().
// after you've implemented snapshots, pass the current snapshot
// (or nil if there's not yet a snapshot).
func (rf *Raft) persist() {
	// Your code here (3C).
	// Example:
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	e.Encode(rf.logs)
	raftstate := w.Bytes()
	rf.persister.Save(raftstate, nil)
}

// restore previously persisted state.
func (rf *Raft) readPersist(data []byte) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if len(data) < 1 { // bootstrap without any state?
		rf.currentTerm = 0
		rf.votedFor = -1
		rf.logs = []LogEntries{}
		rf.state = Follower
		rf.leaderId = -1 // no leader yet
		rf.persist()     // persist the initial state
		return
	}
	r := bytes.NewBuffer(data)
	d := labgob.NewDecoder(r)
	err := d.Decode(rf.currentTerm)
	if err != nil {
		panic(errors.New("could not decode current term from disk"))
	}
	err = d.Decode(rf.votedFor)
	if err != nil {
		panic(errors.New("could not decode votedFor from disk"))
	}

	err = d.Decode(rf.logs)
	if err != nil {
		panic(errors.New("could not logs from disk"))
	}
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (3D).

}

// example RequestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	// Your data here (3A, 3B).
	Term         int //candidate’s term
	CandidateId  int // candidate requesting vote
	LastLogIndex int // index of candidate’s last log entry (§5.4)
	LastLogTerm  int // term of candidate’s last log entry (§5.4)
}

// example RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	// Your data here (3A).
	Term        int  // currentTerm, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (3B).

	return index, term, isLeader
}

// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

func (rf *Raft) AppendEntry(args *AppendEntry, reply *AppendEntryReply) {
	if !rf.killed() {
		rf.mu.Lock()
		defer rf.mu.Unlock()
		// check that this leader is not stale
		if args.Term < rf.currentTerm {
			reply.Term = rf.currentTerm
			reply.Success = false
			return
		}
		// set ourself as follower
		rf.state = Follower
		rf.currentTerm = args.Term
		// persist state
		rf.persist()
		// reset our election timeout; in this case just set the last heartbeat time
		rf.timeout = time.Now().Add(time.Duration(500+rand.Int63()%200) * time.Millisecond)
		rf.leaderId = args.LeaderId

		reply.Term = rf.currentTerm
		reply.Success = true
	}
}

func (rf *Raft) ticker() {
	for !rf.killed() {
		rf.mu.Lock()
		// Check if a leader election should be started.
		if rf.state != Leader && time.Now().After(rf.timeout) {
			// We haven't received a heartbeat in a while, so we should start an election.
			rf.mu.Unlock()
			// start an election
			rf.startElection()

			// Re-lock to update timeout AFTER election starts
			rf.mu.Lock()
			ms := 500 + (rand.Int63() % 200)
			rf.timeout = time.Now().Add(time.Duration(ms) * time.Millisecond)
			rf.mu.Unlock()
		} else {
			rf.mu.Unlock()
		}

		// pause for a random amount of time between 50 and 350
		// milliseconds.
		ms := 50 + (rand.Int63() % 300)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

func (rf *Raft) startElection() {
	// we haven't received an heartbeat in a while, let's start an election
	// 1.	Increment currentTerm and switch to candidate
	// 2.	Vote for self
	// 3.	Reset election timeout
	// 4.	Send RequestVote RPCs to all other servers
	rf.mu.Lock()
	defer rf.mu.Unlock()
	// increment term and become candidate
	rf.currentTerm = rf.currentTerm + 1
	rf.state = Candidate
	rf.votedFor = rf.me
	rf.persist() // save term and vote
	DPrintf("Server %d starting election for term %d due to timeout %v %v", rf.me, rf.currentTerm, time.Now(), rf.timeout)
	var voteMu sync.Mutex
	votes := 1 // we have one vote for ourself
	for peer, _ := range rf.peers {
		if peer == rf.me {
			continue
		}
		lastLogTerm := 0
		lastLogIndex := len(rf.logs) - 1
		if lastLogIndex >= 0 {
			lastLogTerm = rf.logs[lastLogIndex].Term
		}

		args := &RequestVoteArgs{
			Term:         rf.currentTerm,
			CandidateId:  rf.me,
			LastLogIndex: lastLogIndex,
			LastLogTerm:  lastLogTerm,
		}
		go func(peer int) {
			reply := RequestVoteReply{}
			ok := rf.sendRequestVote(peer, args, &reply)
			if ok {
				rf.mu.Lock()
				if reply.Term > rf.currentTerm {
					// we are stale, revert to follower
					DPrintf("Server %d received a higher term %d from %d, reverting to follower", rf.me, reply.Term, peer)
					rf.currentTerm = reply.Term
					rf.votedFor = -1
					rf.persist()
					// set our state to follower
					rf.state = Follower
					rf.mu.Unlock()
					return
				}
				// check if we have won the election but if we are follower it would not matter
				if reply.VoteGranted && rf.state == Candidate {
					rf.mu.Unlock()
					voteMu.Lock()
					votes += 1
					if votes > len(rf.peers)/2 {
						// we have won the election
						DPrintf("Server %d won the election for term %d with %d votes", rf.me, rf.currentTerm, votes)
						voteMu.Unlock()
						rf.becomeLeader()
					} else {
						voteMu.Unlock()
					}
				}
			}
		}(peer)
	}
}

func (rf *Raft) becomeLeader() {
	rf.mu.Lock()
	rf.state = Leader
	rf.leaderId = rf.me // set ourself as leader
	// start sending heartbeats
	term := rf.currentTerm
	rf.mu.Unlock()
	for !rf.killed() {
		rf.mu.Lock()
		if rf.state != Leader {
			// if we are not leader anymore, stop sending heartbeats
			DPrintf("Server %d stopped sending heartbeats because it is no longer leader", rf.me)
			rf.mu.Unlock()
			break
		}
		DPrintf("Server %d sending heartbeats for term %d", rf.me, rf.currentTerm)
		// send heartbeats to all peers
		args := &AppendEntry{
			Term:     term,
			LeaderId: rf.me,
		}
		for peer, _ := range rf.peers {
			reply := &AppendEntryReply{}
			go rf.sendHeartBeat(peer, args, reply)
		}
		rf.mu.Unlock()
		// sleep this thread for 100 milliseconds then send append entries again
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

func (rf *Raft) sendHeartBeat(peer int, args *AppendEntry, reply *AppendEntryReply) {
	ok := rf.peers[peer].Call("Raft.AppendEntry", args, reply)
	if ok {
		rf.mu.Lock()
		defer rf.mu.Unlock()
		if reply.Term > rf.currentTerm {
			// I am stale, this server has a more up-to-date log and hence some other leader with up-to-date log
			rf.state = Follower
			rf.currentTerm = args.Term
			rf.votedFor = -1
			rf.persist()
			rf.leaderId = args.LeaderId
		}
	}
}

// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	if !rf.killed() {
		DPrintf("Server %d received RequestVote from %d: %v", rf.me, args.CandidateId, args)
		rf.mu.Lock()
		defer rf.mu.Unlock()
		// Your code here (3A, 3B).
		if args.Term < rf.currentTerm {
			// this is a stale worker, update him; or we already voted
			DPrintf("Server %d rejected RequestVote from %d: %v", rf.me, args.CandidateId, args)
			reply.Term = rf.currentTerm
			reply.VoteGranted = false
			return
		}
		if args.Term > rf.currentTerm {
			// update our term and step down if necessary
			rf.currentTerm = args.Term
			rf.votedFor = -1
			rf.state = Follower
			rf.persist()
		}
		reply.Term = rf.currentTerm

		// Check if we can vote for this candidate
		if (rf.votedFor == -1 || rf.votedFor == args.CandidateId) && candidateUpToDate(rf, args) {
			rf.votedFor = args.CandidateId
			rf.persist()
			reply.VoteGranted = true
			DPrintf("Server %d voted for %d in term %d", rf.me, args.CandidateId, rf.currentTerm)
		} else {
			reply.VoteGranted = false
			DPrintf("Server %d rejected vote for %d in term %d", rf.me, args.CandidateId, rf.currentTerm)
		}
	}
}

func candidateUpToDate(rf *Raft, args *RequestVoteArgs) bool {
	lastLogIndex := len(rf.logs) - 1
	lastLogTerm := 0
	if lastLogIndex >= 0 {
		lastLogTerm = rf.logs[lastLogIndex].Term
	}

	if args.LastLogTerm != lastLogTerm {
		return args.LastLogTerm > lastLogTerm
	}
	return args.LastLogIndex >= lastLogIndex
}

// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.state = Follower
	rf.timeout = time.Now().Add(time.Duration(500+rand.Int63()%200) * time.Millisecond)

	// Your initialization code here (3A, 3B, 3C).
	// initialize from state persisted before a crash
	go rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()

	return rf
}
