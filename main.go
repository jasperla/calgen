package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func calgen(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("calgen.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		pdf(w, r)
	}
}

func style(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "style.css")
}

func pdf(w http.ResponseWriter, r *http.Request) {
	weeks, _ := strconv.ParseInt(r.Form["weeks"][0], 10, 0)

	// Turn our HTML RFC3339 date into a proper time.Time
	dateString := fmt.Sprintf("%sT00:00:00+00:00", r.Form["begindate"][0])
	dateTime, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		t, _ := template.New("failedTimeParse").Parse("<script>alert('Failed to parse date; please go back')</script>")
		t.Execute(w, nil)

		retry, _ := template.ParseFiles("calgen.gtpl")
		retry.Execute(w, nil)
		return
	}

	f, err := os.Create("output.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	c := bufio.NewWriter(f)

	// Start by sending the header which is actually static.
	t, _ := template.ParseFiles("header.gtpl")
	t.Execute(c, "T")

	// Get the weekday index for the start date
	startDayIdx := int(dateTime.Weekday())

	// Now generate as many empty cells as needed and print the actual
	//    start date
	for i := 0; i < startDayIdx; i++ {
		fmt.Fprintf(c, "<td></td>")
	}
	fmt.Fprint(c, fmt.Sprintf("<td>%02d-%02d</td>", dateTime.Day(), dateTime.Month()))

	// And pad the rest of the week
	for i := 0; i < (7 - startDayIdx - 1); i++ {
		nDay := time.Duration(i*24) * time.Hour
		nextDay := dateTime.Add(time.Hour*24 + nDay)
		fmt.Fprintf(c, fmt.Sprintf("<td>%02d-%02d</td>", nextDay.Day(), nextDay.Month()))
	}

	// For the remaining weeks just fill them out
	for week := 1; week < int(weeks); week++ {
		fmt.Fprintf(c, "<tr>")
		for i := 0; i < 7; i++ {
			offset := time.Duration(i*24) * time.Hour
			// Account for the week offset
			offset += ((time.Hour * 24) * 7) * time.Duration(week*1)
			// Add the week offset, and deduct the first day
			today := dateTime.Add(time.Hour*24 + offset - (time.Hour * 24))
			fmt.Fprintf(c, fmt.Sprintf("<td>%02d-%02d</td>", today.Day(), today.Month()))
		}
		fmt.Fprintf(c, "</tr>")
	}

	// Spit out the footer and be done.
	fmt.Fprintf(c, "</table></body></html>")
	c.Flush()

	// wkhtmltopdf needs a display to run.
	var cmd *exec.Cmd
	if os.Getenv("DISPLAY") == "" {
		cmd = exec.Command("xvfb-run", "wkhtmltopdf", "-O", "landscape", "output.html", "output.pdf")
	} else {
		cmd = exec.Command("wkhtmltopdf", "-O", "landscape", "output.html", "output.pdf")
	}

	out, err := cmd.Output()

	if err != nil {
		fmt.Println("wkhtmltopdf failed:")
		fmt.Println(string(out))
		fmt.Println(err)
		return
	}

	fmt.Print(string(out))

	streamPDFbytes, err := ioutil.ReadFile("./output.pdf")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b := bytes.NewBuffer(streamPDFbytes)

	w.Header().Set("Content-type", "application/pdf")

	if _, err := b.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	w.Write([]byte("PDF Generated"))

	os.Remove("output.html")
	os.Remove("output.pdf")
}

func main() {
	http.HandleFunc("/", calgen)
	http.HandleFunc("/style.css", style)
	http.HandleFunc("/pdf", pdf)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
