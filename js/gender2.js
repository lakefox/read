/* Gender */

// function getWeights(target) {
//   let keys = "abcdefghijklmnopqrstuvwxyz".split("");
//   let letters = {};
//   for (let i = 0; i < keys.length; i++) {
//     letters[keys[i]] = 0;
//   }

//   // Weigh Target names
//   for (let i = 0; i < target.length; i++) {
//     for (let b = 0; b < target[i].length; b++) {
//       letters[target[i][b]] += 1;
//     }
//   }
//   return Object.keys(letters).sort((a, b) => {
//     return letters[b] - letters[a];
//   });
// }

// src https://www.ssa.gov/OACT/babynames/decades/century.html
// let rows = document.querySelectorAll("tr[align='right']");

// let male = [...rows].map((e) => {
//   return e.children[1].innerText.toLowerCase().split("");
// });
// let female = [...rows].map((e) => {
//   return e.children[3].innerText.toLowerCase().split("");
// });

// let maleWeights = getWeights(male);

// let femaleWeights = getWeights(female);

// /* Race */

// let white = [...document.querySelectorAll("td > a")].map((e) => {
//   return e.innerText.toLowerCase();
// });

// let whiteWeights = getWeights(white);

/* Age */
// https://www.babycenter.com/baby-names/most-popular/top-baby-names-2000

// let names = [...document.querySelectorAll("tr > td > a")].map((e) =>
//   e.innerText.toLowerCase()
// );
// const column1 = [];
// const column2 = [];

// for (let i = 0; i < names.length; i += 2) {
//   column1.push(names[i]);
//   column2.push(names[i + 1]);
// }

// // Female
// console.log(getWeights(column1));
// // Male
// console.log(getWeights(column2));

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
const whiteLetters = [
  "e",
  "r",
  "n",
  "a",
  "o",
  "l",
  "s",
  "i",
  "t",
  "h",
  "c",
  "d",
  "m",
  "b",
  "g",
  "u",
  "y",
  "k",
  "w",
  "p",
  "f",
  "v",
  "z",
  "j",
  "x",
  "q",
];
const hispanicLetters = [
  "a",
  "e",
  "o",
  "r",
  "n",
  "l",
  "i",
  "s",
  "c",
  "d",
  "t",
  "m",
  "u",
  "g",
  "z",
  "b",
  "v",
  "p",
  "h",
  "y",
  "j",
  "f",
  "w",
  "q",
  "k",
  "x",
];
const blackLetters = [
  "e",
  "r",
  "a",
  "n",
  "o",
  "l",
  "s",
  "i",
  "t",
  "d",
  "c",
  "m",
  "h",
  "b",
  "y",
  "g",
  "u",
  "w",
  "k",
  "p",
  "f",
  "v",
  "j",
  "z",
  "x",
  "q",
];

const decadeLettersMale = {
  //   2020: [
  //     "a",
  //     "n",
  //     "e",
  //     "o",
  //     "i",
  //     "l",
  //     "r",
  //     "h",
  //     "s",
  //     "m",
  //     "c",
  //     "t",
  //     "d",
  //     "j",
  //     "u",
  //     "y",
  //     "b",
  //     "k",
  //     "v",
  //     "w",
  //     "g",
  //     "x",
  //     "z",
  //     "p",
  //     "f",
  //     "q",
  //   ],
  2010: [
    "a",
    "n",
    "e",
    "i",
    "o",
    "r",
    "l",
    "h",
    "s",
    "c",
    "d",
    "t",
    "j",
    "m",
    "y",
    "b",
    "u",
    "v",
    "k",
    "w",
    "g",
    "p",
    "x",
    "z",
    "f",
    "q",
  ],
  2000: [
    "a",
    "n",
    "e",
    "r",
    "i",
    "o",
    "s",
    "t",
    "l",
    "h",
    "c",
    "j",
    "d",
    "m",
    "u",
    "b",
    "k",
    "y",
    "v",
    "g",
    "p",
    "w",
    "x",
    "z",
    "f",
    "q",
  ],
  1990: [
    "a",
    "e",
    "n",
    "r",
    "t",
    "i",
    "o",
    "s",
    "l",
    "h",
    "d",
    "c",
    "j",
    "y",
    "m",
    "b",
    "u",
    "k",
    "v",
    "p",
    "g",
    "w",
    "f",
    "x",
    "z",
    "q",
  ],
  1980: [
    "a",
    "e",
    "r",
    "n",
    "s",
    "i",
    "l",
    "o",
    "t",
    "h",
    "d",
    "y",
    "c",
    "j",
    "m",
    "u",
    "b",
    "k",
    "p",
    "g",
    "w",
    "f",
    "v",
    "x",
    "z",
    "q",
  ],
  1970: [
    "r",
    "a",
    "e",
    "n",
    "o",
    "l",
    "t",
    "i",
    "d",
    "y",
    "s",
    "h",
    "c",
    "j",
    "m",
    "b",
    "p",
    "g",
    "k",
    "u",
    "w",
    "f",
    "v",
    "q",
    "x",
    "z",
  ],
  1960: [
    "r",
    "e",
    "a",
    "n",
    "l",
    "i",
    "t",
    "y",
    "o",
    "d",
    "h",
    "c",
    "m",
    "s",
    "j",
    "b",
    "g",
    "p",
    "k",
    "w",
    "f",
    "u",
    "v",
    "q",
    "x",
    "z",
  ],
  1950: [
    "e",
    "r",
    "a",
    "n",
    "l",
    "i",
    "o",
    "d",
    "t",
    "h",
    "y",
    "s",
    "c",
    "m",
    "g",
    "p",
    "b",
    "f",
    "j",
    "u",
    "k",
    "w",
    "v",
    "q",
    "x",
    "z",
  ],
  1940: [
    "e",
    "r",
    "a",
    "l",
    "n",
    "o",
    "i",
    "d",
    "h",
    "y",
    "m",
    "t",
    "s",
    "c",
    "b",
    "j",
    "g",
    "p",
    "u",
    "w",
    "f",
    "k",
    "v",
    "q",
    "x",
    "z",
  ],
};

const decadeLettersFemale = {
  //   2020: [
  //     "a",
  //     "l",
  //     "e",
  //     "i",
  //     "n",
  //     "y",
  //     "r",
  //     "o",
  //     "m",
  //     "h",
  //     "s",
  //     "c",
  //     "v",
  //     "d",
  //     "t",
  //     "b",
  //     "k",
  //     "p",
  //     "g",
  //     "u",
  //     "z",
  //     "w",
  //     "j",
  //     "f",
  //     "q",
  //     "x",
  //   ],
  2010: [
    "a",
    "l",
    "e",
    "n",
    "i",
    "y",
    "r",
    "o",
    "m",
    "s",
    "h",
    "b",
    "c",
    "t",
    "d",
    "k",
    "g",
    "v",
    "u",
    "j",
    "p",
    "x",
    "z",
    "f",
    "q",
    "w",
  ],
  2000: [
    "a",
    "e",
    "i",
    "n",
    "l",
    "r",
    "s",
    "y",
    "h",
    "m",
    "t",
    "c",
    "b",
    "o",
    "d",
    "k",
    "g",
    "j",
    "u",
    "v",
    "x",
    "p",
    "z",
    "f",
    "q",
    "w",
  ],
  1990: [
    "a",
    "e",
    "i",
    "n",
    "l",
    "r",
    "s",
    "t",
    "h",
    "y",
    "c",
    "m",
    "k",
    "o",
    "b",
    "d",
    "j",
    "g",
    "u",
    "p",
    "v",
    "f",
    "x",
    "q",
    "w",
    "z",
  ],
  1980: [
    "a",
    "e",
    "i",
    "n",
    "r",
    "l",
    "s",
    "t",
    "c",
    "y",
    "h",
    "m",
    "d",
    "k",
    "o",
    "u",
    "b",
    "j",
    "p",
    "f",
    "g",
    "v",
    "q",
    "w",
    "z",
    "x",
  ],
  1970: [
    "a",
    "e",
    "n",
    "i",
    "r",
    "l",
    "t",
    "c",
    "h",
    "s",
    "y",
    "d",
    "m",
    "o",
    "b",
    "k",
    "j",
    "u",
    "p",
    "g",
    "f",
    "v",
    "w",
    "z",
    "q",
    "x",
  ],
  1960: [
    "a",
    "e",
    "n",
    "i",
    "r",
    "l",
    "t",
    "h",
    "y",
    "c",
    "d",
    "o",
    "s",
    "b",
    "j",
    "m",
    "k",
    "u",
    "g",
    "p",
    "v",
    "w",
    "z",
    "f",
    "q",
    "x",
  ],
  1950: [
    "a",
    "e",
    "n",
    "r",
    "i",
    "l",
    "t",
    "h",
    "s",
    "c",
    "o",
    "y",
    "d",
    "j",
    "b",
    "m",
    "u",
    "g",
    "k",
    "v",
    "p",
    "z",
    "f",
    "q",
    "w",
    "x",
  ],
  1940: [
    "a",
    "e",
    "n",
    "r",
    "l",
    "i",
    "o",
    "t",
    "s",
    "y",
    "d",
    "h",
    "j",
    "c",
    "m",
    "u",
    "b",
    "g",
    "p",
    "k",
    "v",
    "f",
    "w",
    "q",
    "z",
    "x",
  ],
};

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

function estimateGender(name) {
  const maleScore = calculateScore(name, maleLetters);
  const femaleScore = calculateScore(name, femaleLetters);

  if (maleScore > femaleScore) {
    return "male";
  } else if (femaleScore >= maleScore) {
    return "female";
  }
}

function estimateRace(name) {
  const blackScore = calculateScore(name, blackLetters);
  const hispanicScore = calculateScore(name, hispanicLetters);
  const whiteScore = calculateScore(name, whiteLetters);

  let res = {
    black: blackScore,
    hispanic: hispanicScore,
    white: whiteScore,
  };

  let scores = Object.keys(res).sort((a, b) => {
    return res[b] - res[a];
  });
  return scores[0];
}

function estimateAge(name, gender) {
  let res = {};
  if (gender == "male") {
    for (const decade in decadeLettersMale) {
      res[decade] = calculateScore(name, decadeLettersMale[decade]);
    }
  } else {
    for (const decade in decadeLettersFemale) {
      res[decade] = calculateScore(name, decadeLettersFemale[decade]);
    }
  }

  let scores = Object.keys(res).sort((a, b) => {
    return res[b] - res[a];
  });
  return [scores[0], scores[1]];
}

// Example usage

function profile(name) {
  let first = name.toLowerCase().split(" ")[0];
  let last = name.toLowerCase().split(" ")[1];
  let gender = estimateGender(first);
  let race = estimateRace(last);
  let age = estimateAge(first, gender);
  age = age
    .sort((a, b) => {
      return parseInt(a) - parseInt(b);
    })
    .join(" to ");
  return { gender, age, race };
}
