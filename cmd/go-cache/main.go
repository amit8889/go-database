package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amit8889/go-cache/internal/config"
)

type Node struct {
	data  string
	left  *Node
	right *Node
}
type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}
type Cache struct {
	Queue Queue
	Hash  Hash
}
type Hash map[string]*Node

func NewCache() Cache {
	return Cache{Queue: NewQueue(), Hash: Hash{}}
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}
	head.right = tail
	head.left = nil
	tail.left = head
	tail.data = ""
	tail.right = nil
	return Queue{Head: head, Tail: tail, Length: 0}
}

func (c *Cache) Check(str string) {
	var node *Node
	if val, ok := c.Hash[str]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{data: str}
	}
	c.Add(node)
	c.Hash[str] = node
}

func (c *Cache) Add(node *Node) {
	if c.Queue.Length == 0 {
		c.Queue.Head.right = node
		node.left = c.Queue.Head
		c.Queue.Tail.left = node
		node.right = c.Queue.Tail
		c.Queue.Length++
		return
	}
	if c.Queue.Length == 1 {
		node.left = c.Queue.Head
		node.right = c.Queue.Tail
		c.Queue.Head.right = node
		c.Queue.Tail.left = node
		c.Queue.Length++
		return
	}

	node.left = c.Queue.Head.right
	node.right = c.Queue.Head.right.right
	c.Queue.Head.right.right.left = node
	c.Queue.Head.right.right = node
	c.Queue.Length++
}

func (c *Cache) Remove(node *Node) *Node {
	if node.left != nil {
		node.left.right = node.right
	}
	if node.right != nil {
		node.right.left = node.left
	}
	node.left = nil
	node.right = nil
	c.Queue.Length--
	return node
}

func (c *Cache) Display() {
	current := c.Queue.Head.right
	for current != c.Queue.Tail {
		fmt.Printf("%s -> ", current.data)
		current = current.right
	}
	fmt.Println("nil")
}

func main() {
	fmt.Println("main function is called")
	cfg := config.MustLoad()
	fmt.Println("config loaded successfully")
	// cache setup
	cache := NewCache()
	for _, word := range []string{"node", "go lang", "react", "java", "rabbit mq"} {
		cache.Check(word)
		cache.Display()
	}

	// server setup
	server := http.Server{
		Addr: cfg.Server.Host + ":" + cfg.Server.PORT,
	}
	fmt.Println("server is running on port", cfg.Server.PORT)

	// make a go channel for signal handling
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-done
	log.Println("server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Error shutting down server:", err)
	}
	log.Println("server is shut down")
}
