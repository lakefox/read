// {
//               title: "Audio Title",
//               artist: "Artist Name",
//               album: "Album Name",
//               artwork: [
//                 {
//                   src: "https://media.wired.com/photos/666c7786c0e3c3ad99f510cc/191:100/w_1280,c_limit/laurai%20-%208.JPG",
//                   sizes: "512x512",
//                   type: "image/jpeg",
//                 },
//               ],
//             }

export function streamAudio(textInput, audioPlayer, trackInfo) {
  return new Promise(async (resolve, reject) => {
    try {
      const response = await fetch("https://143.244.148.224/stream", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ data: cleanTextForCLI(textInput) }),
      });

      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      let mediaSource;
      if (window.MediaSource != undefined) {
        mediaSource = new MediaSource();
      } else {
        mediaSource = new ManagedMediaSource();
      }

      audioPlayer.src = URL.createObjectURL(mediaSource);
      console.log(audioPlayer.src);

      mediaSource.addEventListener("sourceopen", async () => {
        const sourceBuffer = mediaSource.addSourceBuffer(
          'audio/webm; codecs="opus"'
        );

        const reader = response.body.getReader();
        let streamQueue = [];
        let maxBufferSize = 6000; // Maximum buffer size in seconds

        const processChunk = async ({ done, value }) => {
          if (done) {
            mediaSource.endOfStream();
            return;
          }
          streamQueue.push(value);
          if (!sourceBuffer.updating) {
            appendBuffer();
          }
        };

        const appendBuffer = () => {
          if (sourceBuffer.updating || streamQueue.length === 0) {
            return;
          }

          // Remove old buffered data if buffer is full
          if (sourceBuffer.buffered.length > 0) {
            const bufferedEnd = sourceBuffer.buffered.end(
              sourceBuffer.buffered.length - 1
            );
            const bufferedStart = sourceBuffer.buffered.start(0);
            if (bufferedEnd - bufferedStart > maxBufferSize) {
              sourceBuffer.remove(bufferedStart, bufferedStart + 5);
              return;
            }
          }

          sourceBuffer.appendBuffer(streamQueue.shift());
        };

        sourceBuffer.addEventListener("updateend", appendBuffer);

        while (true) {
          const chunk = await reader.read();
          await processChunk(chunk);
        }
      });

      // Set up Media Session API
      if ("mediaSession" in navigator) {
        navigator.mediaSession.metadata = new MediaMetadata(trackInfo);

        // Set action handlers (optional)
        navigator.mediaSession.setActionHandler("play", () => {
          audioPlayer.play();
        });
        navigator.mediaSession.setActionHandler("pause", () => {
          audioPlayer.pause();
        });
        navigator.mediaSession.setActionHandler("seekbackward", (details) => {
          audioPlayer.currentTime = Math.max(
            audioPlayer.currentTime - (details.seekOffset || 10),
            0
          );
        });
        navigator.mediaSession.setActionHandler("seekforward", (details) => {
          audioPlayer.currentTime = Math.min(
            audioPlayer.currentTime + (details.seekOffset || 10),
            audioPlayer.duration
          );
        });
        navigator.mediaSession.setActionHandler("previoustrack", () => {
          // Handle previous track action
        });
        navigator.mediaSession.setActionHandler("nexttrack", () => {
          // Handle next track action
        });
      }

      resolve(audioPlayer);
    } catch (error) {
      console.error("Error:", error);
    }
  });
}
export function cleanTextForCLI(text) {
  // Replace potentially harmful characters with safe ones
  let sanitizedText = text.replace(/[`$&|;<>]/g, "");

  // Escape single quotes, double quotes, and backslashes
  sanitizedText = sanitizedText.replace(/[^a-zA-Z0-9\s\.]/g, "");

  // Trim whitespace
  sanitizedText = sanitizedText.trim();

  // sanitizedText = sanitizedText
  //   .replace(/[^a-zA-Z0-9\s\.]/g, " ")
  //   .replace(/\s+/g, " ");

  return sanitizedText;
}

export function textToSpeech(text, options = {}, callback = null) {
  // Create a new SpeechSynthesisUtterance instance
  let utterance = new SpeechSynthesisUtterance(text);

  // Set default options
  let defaultOptions = {
    lang: "en-US",
    voice: null, // Set to a specific voice if needed
    pitch: 1,
    rate: 0.8,
    volume: 1,
  };

  // Override default options with user-provided options
  let config = { ...defaultOptions, ...options };

  // Set the properties on the utterance instance
  utterance.lang = config.lang;
  utterance.pitch = config.pitch;
  utterance.rate = config.rate;
  utterance.volume = config.volume;

  // Find and set the desired voice if provided
  if (config.voice) {
    let voices = window.speechSynthesis.getVoices();
    let selectedVoice = voices.find((voice) => voice.name === config.voice);
    if (selectedVoice) {
      utterance.voice = selectedVoice;
    }
  }

  // Set the onend callback if provided
  if (callback) {
    utterance.onend = callback;
  }

  // Speak the text
  window.speechSynthesis.speak(utterance);
}

// Example usage:
// textToSpeech(
//   "Hello, world!",
//   {
//     lang: "en-US",
//     voice: "Google UK English Male",
//     pitch: 1.2,
//     rate: 1,
//     volume: 0.8,
//   },
//   () => {
//     console.log("Speech synthesis finished.");
//   }
// );
