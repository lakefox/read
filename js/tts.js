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
