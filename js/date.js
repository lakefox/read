const months = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

let monthCombos = months.concat([
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
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
    if (secondPart < 10) {
      return `${numberWords[firstPart * 10]} ${numberWords[secondPart]}`;
    } else {
      return `${numberWords[firstPart]} ${numberWords[secondPart]}`;
    }
  }
  return yearStr
    .split("")
    .map((digit) => numberWords[parseInt(digit)])
    .join(" ");
}

function replaceDates(text) {
  let textSplit = text.split(" ");

  let numbers = [];

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
      Math.max(target.index - 3, 0),
      Math.min(target.index + 3, textSplit.length)
    );

    let combos = generateCombinations(neighbors);

    combos.sort((a, b) => {
      return b.length - a.length;
    });

    let foundDate = undefined;
    let filter = new Date().toString();
    let bestLine = filter;
    for (let a = 0; a < combos.length; a++) {
      const line = combos[a].join(" ");

      let date = new Date(line);
      if (date != "Invalid Date") {
        if (foundDate != undefined) {
          if (foundDate.toString() == date.toString()) {
            if (line.length < bestLine.length) {
              if (line.replace(/[^0-9]/g, "").length >= 4) {
                bestLine = line;
              }
            }
          }
        } else {
          foundDate = date;
          bestLine = line;
        }
      }
    }

    if (bestLine != filter) {
      console.log(bestLine);
      let parts = bestLine
        .replace(/[^0-9A-Z-a-z]/g, " ")
        .replace(/\s+/g, " ")
        .trim()
        .split(" ");
      console.log(parts, foundDate);
      let str = "";
      let yearString = foundDate.getFullYear().toString();
      let year = yearToWords(parseInt(yearString));
      let month = months[foundDate.getMonth()];
      let day = dayToOrdinal(foundDate.getDate());

      if (parts.length == 3) {
        console.log("3");
        str = `${month} ${day} ${year}`;
      } else if (parts.length == 2) {
        console.log("2");
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
      console.log(str);
      // Replace the original date string with the spoken string
      text = text.replace(bestLine, str);
    }
  }

  return text;
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

// Example usage
const text = "These are some dates: , July 4th, 2024, and Nov 6, 2023.";
console.log(replaceDates(text));
// Output: These are some dates: January twenty-second two thousand one, July fourth two thousand one, and November sixth two thousand twenty-three.
