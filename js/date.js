const months = [
  "january",
  "february",
  "march",
  "april",
  "may",
  "june",
  "july",
  "august",
  "september",
  "october",
  "november",
  "december",
];

let monthCombos = months.concat([
  "jan",
  "feb",
  "mar",
  "apr",
  "jun",
  "jul",
  "aug",
  "sep",
  "oct",
  "nov",
  "dec",
]);

const ordinals = [
  "zeroth",
  "first",
  "second",
  "third",
  "fourth",
  "fifth",
  "sixth",
  "seventh",
  "eighth",
  "ninth",
  "tenth",
  "eleventh",
  "twelfth",
  "thirteenth",
  "fourteenth",
  "fifteenth",
  "sixteenth",
  "seventeenth",
  "eighteenth",
  "nineteenth",
  "twentieth",
  "twenty-first",
  "twenty-second",
  "twenty-third",
  "twenty-fourth",
  "twenty-fifth",
  "twenty-sixth",
  "twenty-seventh",
  "twenty-eighth",
  "twenty-ninth",
  "thirtieth",
  "thirty-first",
];

const numberWords = [
  "zero",
  "one",
  "two",
  "three",
  "four",
  "five",
  "six",
  "seven",
  "eight",
  "nine",
  "ten",
  "eleven",
  "twelve",
  "thirteen",
  "fourteen",
  "fifteen",
  "sixteen",
  "seventeen",
  "eighteen",
  "nineteen",
  "twenty",
  "twenty-one",
  "twenty-two",
  "twenty-three",
  "twenty-four",
  "twenty-five",
  "twenty-six",
  "twenty-seven",
  "twenty-eight",
  "twenty-nine",
  "thirty",
  "thirty-one",
  "thirty-two",
  "thirty-three",
  "thirty-four",
  "thirty-five",
  "thirty-six",
  "thirty-seven",
  "thirty-eight",
  "thirty-nine",
  "forty",
  "forty-one",
  "forty-two",
  "forty-three",
  "forty-four",
  "forty-five",
  "forty-six",
  "forty-seven",
  "forty-eight",
  "forty-nine",
  "fifty",
  "fifty-one",
  "fifty-two",
  "fifty-three",
  "fifty-four",
  "fifty-five",
  "fifty-six",
  "fifty-seven",
  "fifty-eight",
  "fifty-nine",
  "sixty",
  "sixty-one",
  "sixty-two",
  "sixty-three",
  "sixty-four",
  "sixty-five",
  "sixty-six",
  "sixty-seven",
  "sixty-eight",
  "sixty-nine",
  "seventy",
  "seventy-one",
  "seventy-two",
  "seventy-three",
  "seventy-four",
  "seventy-five",
  "seventy-six",
  "seventy-seven",
  "seventy-eight",
  "seventy-nine",
  "eighty",
  "eighty-one",
  "eighty-two",
  "eighty-three",
  "eighty-four",
  "eighty-five",
  "eighty-six",
  "eighty-seven",
  "eighty-eight",
  "eighty-nine",
  "ninety",
  "ninety-one",
  "ninety-two",
  "ninety-three",
  "ninety-four",
  "ninety-five",
  "ninety-six",
  "ninety-seven",
  "ninety-eight",
  "ninety-nine",
];

const yearNumeberWords = [
  "one hundred",
  "two hundred",
  "three hundred",
  "four hundred",
  "five hundred",
  "six hundred",
  "seven hundred",
  "eight hundred",
  "nine hundred",
  "one thousand",
  "eleven",
  "twelve",
  "thirteen",
  "fourteen",
  "fifteen",
  "sixteen",
  "seventeen",
  "eighteen",
  "nineteen",
  "two thousand",
];

function numberToWords(number) {
  if (number <= 99) {
    return numberWords[number];
  }
  return number
    .toString()
    .split("")
    .map((digit) => numberWords[parseInt(digit)])
    .join(" ");
}

function dayToOrdinal(day) {
  if (day <= 31) {
    return ordinals[day];
  }
  return numberToWords(day);
}

function yearToWords(year) {
  const yearStr = year.toString();
  if (yearStr.length === 4) {
    const firstPart = parseInt(yearStr.slice(0, 2));
    const secondPart = parseInt(yearStr.slice(2, 4));

    return `${yearNumeberWords[firstPart - 1]} ${numberWords[secondPart]}`;
  }
  return yearStr
    .split("")
    .map((digit) => numberWords[parseInt(digit)])
    .join(" ");
}

export function replaceDates(text) {
  let textSplit = text.split(" ");

  let numbers = [];

  let replaceable = [];

  for (let i = 0; i < textSplit.length; i++) {
    if (!isNaN(parseInt(textSplit[i]))) {
      numbers.push({
        text: textSplit[i],
        index: i,
      });
    }
  }

  for (let i = 0; i < numbers.length; i++) {
    const target = numbers[i];
    let neighbors = textSplit.slice(
      Math.max(target.index - 2, 0),
      Math.min(target.index + 2, textSplit.length)
    );
    if (i > 0) {
      if (target.index - numbers[i - 1].index <= 1) {
        continue;
      }
    }

    let combos = generateCombinations(neighbors);

    combos.sort((a, b) => {
      return b.length - a.length;
    });

    console.log(combos);

    let foundDate = undefined;
    let filter = new Date().toString();
    let bestLine = filter;
    for (let a = 0; a < combos.length; a++) {
      const line = combos[a].join(" ");

      let date = new Date(keepOnlyMonths(line));

      if (date != "Invalid Date") {
        if (foundDate != undefined) {
          if (line.indexOf(target.text) != -1) {
            if (
              keepOnlyMonths(line).trim().split(" ").length >=
              keepOnlyMonths(bestLine).trim().split(" ").length
            ) {
              bestLine = line;
            }
          }
        } else if (line.indexOf(target.text) != -1) {
          foundDate = date;
          bestLine = line;
        }
      }
    }

    if (bestLine != filter) {
      let parts = bestLine
        .replace(/[^0-9A-Z-a-z]/g, " ")
        .replace(/\s+/g, " ")
        .trim()
        .split(" ");
      let str = "";
      let yearString = foundDate.getFullYear().toString();
      let year = yearToWords(parseInt(yearString));
      let month = months[foundDate.getMonth()];
      let day = dayToOrdinal(foundDate.getDate());

      if (parts.length == 3) {
        str = `${month} ${day} ${year}`;
      } else if (parts.length == 2) {
        if (monthCombos.indexOf(parts[0].toLowerCase()) != -1) {
          str = month;
        }
        if (parseInt(parts[1]).toString() == parts[1] && parts[1].length == 4) {
          str += ` ${year}`;
        } else {
          str += ` ${day}`;
        }
      } else if (parts.length == 1 && parts[0].length == 4) {
        str = year;
      }

      replaceable.push({
        text: str,
        index: target.index,
        line: bestLine.split(" "),
      });
    }
  }

  return replaceLinesWithText(replaceable, textSplit).join(" ");
}

function generateCombinations(array) {
  const result = [];

  function helper(start, combination) {
    result.push(combination.slice());

    for (let i = start; i < array.length; i++) {
      combination.push(array[i]);
      helper(i + 1, combination);
      combination.pop();
    }
  }

  helper(0, []);
  return result.slice(1);
}

function replaceLinesWithText(objects, textArray) {
  // Helper function to find the exact start index of the line within the range
  function findLineIndex(line, range) {
    for (let i = range[0]; i <= range[1]; i++) {
      if (
        textArray
          .slice(i, i + line.length)
          .join(" ")
          .toLowerCase()
          .trim() === line.join(" ").toLowerCase().trim()
      ) {
        return i;
      }
    }
    return -1; // If not found
  }

  objects.forEach((obj) => {
    const range = [
      Math.max(0, obj.index - 5),
      Math.min(textArray.length, obj.index + 5),
    ];
    const startIndex = findLineIndex(obj.line, range);

    if (startIndex !== -1) {
      textArray.splice(startIndex, obj.line.length, obj.text);
    }
  });

  return textArray;
}

function keepOnlyMonths(input) {
  const months = [
    "january",
    "jan",
    "february",
    "feb",
    "march",
    "mar",
    "april",
    "apr",
    "may",
    "june",
    "jun",
    "july",
    "jul",
    "august",
    "aug",
    "september",
    "sep",
    "october",
    "oct",
    "november",
    "nov",
    "december",
    "dec",
  ];
  input = input.toLowerCase();
  // Create a regular expression to match all months

  // Initialize an empty string for the result
  let result = "";

  // Split the input by non-alphabetical characters
  const parts = input.split(/([^a-zA-Z]+)/);

  // Iterate through the parts and keep only the months and special characters
  for (let part of parts) {
    if (months.includes(part)) {
      result += part;
    } else if (/[^a-zA-Z]/.test(part)) {
      result += part;
    }
  }

  return result;
}

// Example usage
// const text =
//   "Jackson committed what some would call a betrayal of trust and others would call an act of heroism. Donald Chapman says, “I’ve served my time and now I should be left alone.”Photograph by Eileen TravellOn November 17, 1992, the elderly parents of Donald Arthur Chapman drove to Avenel, New Jersey, to bring their son home. Donald Chapman had been away since 1980, when he turned himself in to police and confessed to kidnapping and raping a young woman. Few of the Chapmans’ neighbors knew the details of his crime. There had been little publicity and no trial, and Chapman’s homecoming failed to stir much interest in the handsome, middle-class town of Wyckoff, where he had grown up and his parents still lived.Avenel is twelve miles south of Newark on Route 1, past a go-go cocktail lounge and a XXX video store. Until recently, a billboard advertising the Hot Tub Club in the Post Road Inn (“Get Wet! $26.95”) marked the turnoff to the Adult Diagnostic and Treatment Center, one of the ";
// console.log(replaceDates(text));
// Output: These are some dates: January twenty-second two thousand one, July fourth two thousand one, and November sixth two thousand twenty-three.

// make this in go and the gender2 stuff
// readablity js in to go
// build account system with sql
// rss feeds
// music inject
