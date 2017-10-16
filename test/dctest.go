package main

import (
	"net/http"
	"time"

	"github.com/TingYunAPM/go"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func redirectOKHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{"status":"success","result":"127.0.0.1:90"}`)
	w.Write(body)
}

func redirect500Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	body := []byte(`null pointer`)
	w.Write(body)
}

func redirectToTimeoutHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{"status":"success","result":"192.168.8.8:88"}`)
	w.Write(body)
}

func redirectInvalidHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{invalid`)
	w.Write(body)
}

func redirectFailHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{"status":"error","result":"db error"}`)
	w.Write(body)
}

func initOKHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{"status":"success","result":{"applicationId":2,"enabled":true,"appSessionKey":"113641","dataSentInterval":60,"apdex_t":500,"config":{"nbs.action_tracer.enabled":true,"nbs.action_tracer.action_threshold":500,"nbs.action_tracer.stack_trace_threshold":500}}}`)
	w.Write(body)
}

func initInvalidHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{invalid`)
	w.Write(body)
}

func initFailHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{"status":"error","result":"db error"}`)
	w.Write(body)
}

func initDisableHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{"status":"success","result":{"applicationId":2,"enabled":false,"appSessionKey":"113641","dataSentInterval":60,"apdex_t":500,"config":{"nbs.agent_enabled":false,"nbs.action_tracer.enabled":true,"nbs.action_tracer.action_threshold":500,"nbs.action_tracer.stack_trace_threshold":500}}}`)
	w.Write(body)
}

func uploadOKHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(8)
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{"status":"success","result":{}}`)
	w.Write(body)
}

func uploadInvalidHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{invalid`)
	w.Write(body)
}

func uploadFailHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	body := []byte(`{"status":"error","result":{"errorCode":1, "errorMessage":"json invalid"}}`)
	w.Write(body)
}

func mainOK() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainRedirectTimeout() {
	tingyun.AppInit("dc-timeout.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainRedirect404() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	//tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainRedirect500() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	http.HandleFunc("/getRedirectHost", redirect500Handler)
	tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainRedirectInvalid() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectInvalidHandler)
	//tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainRedirectFail() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectFailHandler)
	//tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainInitTimeout() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectToTimeoutHandler)
	tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainInit404() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	//tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainInitInvalid() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	tingyun.HandleFunc("/initAgentApp", initInvalidHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainInitFail() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	tingyun.HandleFunc("/initAgentApp", initFailHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func mainUpload404() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/", NotFoundHandler)

	http.ListenAndServe(":90", nil)
}

func mainUploadDropData() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadFailHandler)

	http.ListenAndServe(":90", nil)
}

func mainUploadDCFail() {
	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	tingyun.HandleFunc("/initAgentApp", initOKHandler)
	tingyun.HandleFunc("/upload", uploadInvalidHandler)

	http.ListenAndServe(":90", nil)
}

func mainDisableAgent() {

	tingyun.AppInit("dc-private.json")
	defer tingyun.AppStop()

	tingyun.HandleFunc("/getRedirectHost", redirectOKHandler)
	tingyun.HandleFunc("/initAgentApp", initDisableHandler)
	tingyun.HandleFunc("/upload", uploadOKHandler)

	http.ListenAndServe(":90", nil)
}

func main() {
	//mainOK()

	//mainRedirectTimeout()
	//mainRedirect404()
	//mainRedirect500()
	//mainRedirectInvalid()
	//mainRedirectFail()

	//mainInitTimeout()
	//mainInit404()
	//mainInitInvalid()
	//mainInitFail()

	mainUpload404()
	//mainUploadDropData()
	//mainUploadDCFail()

	//mainDisableAgent()

}
