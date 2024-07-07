import { style, State } from "./html.js";

export class Player extends State {
  constructor(page) {
    let { $ } = super();

    let a = document.querySelector("#audio");
    document.querySelector("#background").src = page.image;

    a.src = `https://api.szn.io/?url=${page.url}`;
    a.onload = () => {
      document.getElementById("play-pause").click();
    };

    let trackInfo = {
      title: page.title,
      artist: page.byline,
      album: page.site,
      artwork: [
        {
          src: page.image,
          sizes: "512x512",
          type: "image/jpeg",
        },
      ],
    };

    // Set up Media Session API
    if ("mediaSession" in navigator) {
      navigator.mediaSession.metadata = new MediaMetadata(trackInfo);

      // Set action handlers (optional)
      navigator.mediaSession.setActionHandler("play", () => {
        a.play();
      });
      navigator.mediaSession.setActionHandler("pause", () => {
        a.pause();
      });
      navigator.mediaSession.setActionHandler("seekbackward", (details) => {
        a.currentTime = Math.max(a.currentTime - (details.seekOffset || 10), 0);
      });
      navigator.mediaSession.setActionHandler("seekforward", (details) => {
        a.currentTime = Math.min(
          a.currentTime + (details.seekOffset || 10),
          a.duration
        );
      });
    }
  }
}

let css = style(/* css */ ``);

async function getAudio(data) {
  try {
    const response = await fetch("http://143.244.148.224:80/stream", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ data: data }),
    });

    if (!response.ok) {
      throw new Error("Network response was not ok");
    }

    const audioBlob = await response.blob();
    const audioUrl = URL.createObjectURL(audioBlob);
    return audioUrl;
  } catch (error) {
    console.error("Error:", error);
  }
}
