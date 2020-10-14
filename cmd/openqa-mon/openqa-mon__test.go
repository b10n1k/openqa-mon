/*
 * Tests for the main display
 */
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Test jobs
var jobs []Job
var instance string

func getAPIRequestJobId(uri string) (int, error) {
	i := strings.LastIndex(uri, "/")
	if i <= 0 {
		return 0, errors.New("Invalid request URI")
	}
	i, err := strconv.Atoi(uri[i+1:])
	if err != nil {
		return 0, err
	}
	return i, nil
}

func wwwAPIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "api")
}
func wwwAPIJob(w http.ResponseWriter, r *http.Request) {
	jobID, err := getAPIRequestJobId(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}
	for _, job := range jobs {
		if job.ID == jobID {
			ret := make(map[string]Job, 0)
			ret["job"] = job
			buf, err := json.Marshal(ret)
			if err != nil {
				panic(err)
			}
			w.Write(buf)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Job not found")
}
func wwwAPIJobOverview(w http.ResponseWriter, r *http.Request) {
	buf, err := json.Marshal(jobs)
	if err != nil {
		panic(err)
	}
	w.Write(buf)
}
func www404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 error")
}

// Start webserver in background
func setupTestServer(addr string) *http.Server {
	srv := &http.Server{Addr: addr}
	http.HandleFunc("/", www404)
	http.HandleFunc("/api/v1", wwwAPIHandler)
	http.HandleFunc("/api/v1/jobs/", wwwAPIJob)
	http.HandleFunc("/api/v1/jobs/overview", wwwAPIJobOverview)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	return srv
}

// Create test jobs
func setupJobs() {
	jobs = make([]Job, 0)
	jobs = append(jobs, Job{ID: 1, Test: "Test1"})
	jobs = append(jobs, Job{ID: 2, Test: "Test2"})
}

// 0 if equal, -1 if jobs1 contains jobs that are not in jobs 2, and 1 if jobs2 contains jobs that are not in jobs1
func jobs_cmp(jobs1 []Job, jobs2 []Job) int {
	if len(jobs1) > len(jobs2) {
		return -1
	}
	if len(jobs1) < len(jobs2) {
		return 1
	}

	// XXX This can be improved!
	for _, ref := range jobs1 {
		found := false
		for _, job := range jobs2 {
			if ref.cmd(job) {
				found = true
				break
			}
		}
		if !found {
			return 1
		}
	}
	return 0
}

// Main test entry point
func TestMain(m *testing.M) {
	setupJobs()
	instance = "127.0.0.1:9091"
	fmt.Printf("Setting up test server (%s) ... \n", instance)
	srv := setupTestServer(instance)
	time.Sleep(100 * time.Millisecond) // Give the server some time to come up

	// Run tests
	ret := m.Run()
	//time.Sleep(100 * time.Second)

	srv.Shutdown(context.TODO())
	os.Exit(ret)
}

func TestJobOverview(t *testing.T) {
	refJobs := jobs
	link := "http://" + instance
	jobs, err := getJobsOverview(link)
	if err != nil {
		t.Fatalf("Error fetching jobs: %s", err)
		return
	}
	if jobs_cmp(jobs, refJobs) != 0 {
		t.Error("Fetched jobs does not match expected jobs")
		return
	}
}
