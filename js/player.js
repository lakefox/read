import { style, State } from "./html.js";
import { streamAudio, cleanTextForCLI } from "./tts.js";

export class Player extends State {
  constructor(page) {
    let { $ } = super();

    let a = document.querySelector("#audio");
    document.querySelector("#background").src = page.image;

    streamAudio({ text: page.text, voice: page.voice }, a, {
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
    }).then((play) => {
      document.getElementById("play-pause").click();
    });
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
