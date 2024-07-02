const maleLetters = [
  "a",
  "e",
  "r",
  "n",
  "l",
  "o",
  "t",
  "i",
  "h",
  "s",
  "y",
  "d",
  "j",
  "c",
  "b",
  "m",
  "u",
  "g",
  "k",
  "w",
  "p",
  "v",
  "f",
  "x",
  "z",
  "q",
];
const femaleLetters = [
  "a",
  "e",
  "i",
  "n",
  "r",
  "l",
  "h",
  "t",
  "s",
  "c",
  "y",
  "m",
  "o",
  "d",
  "b",
  "j",
  "u",
  "g",
  "k",
  "v",
  "f",
  "p",
  "q",
  "x",
  "z",
  "w",
];

function calculateScore(name, letters) {
  let score = 0;
  name
    .toLowerCase()
    .split("")
    .forEach((char) => {
      const index = letters.indexOf(char);
      if (index !== -1) {
        score += 26 - index;
      }
    });
  return score;
}

export function estimateGender(name) {
  const maleScore = calculateScore(name, maleLetters);
  const femaleScore = calculateScore(name, femaleLetters);

  if (maleScore > femaleScore) {
    return "male";
  } else if (femaleScore >= maleScore) {
    return "female";
  }
}
