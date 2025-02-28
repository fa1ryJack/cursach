# cursach 

The projet I made as assignment for college.

"cursach" - "cur" is read the same as in**cur**, "sach" as such.

# What does it do

In this app you paste SoundCloud user's link and then you can see their last 50 likes tracks grouped by uploaders in a form of interactive [circle packing](https://en.wikipedia.org/wiki/Circle_packing).

You can click on circles to zoom in. Click on zoomed track for SoundCloud widget appear.

# How to run

You need installed `go` and `npm`.

Clone or download this repo.

In `server` directory run `go run .` in CLI.

In parentt diretory run `npm install` and `npm run dev`.

# What it uses

The server is written with go and it only fetches likes from SoundCloud.

I use [this](https://github.com/zackradisic/soundcloud-api/tree/master) wrapper for SoundCloud API, although I had to rewrite algorithm for getting client_id.

The front is written with Vue 3 and uses D3.js to visualise likes. I forked an example of circle packing and tweaked with it, you can check it out [here](https://observablehq.com/d/14842d00f4787ffb) (you can also treat it as a demo).

# TODO

- choose how many likes?
  
I think this app has some potential and could be improved but I'm not interested in it (at least currently).
