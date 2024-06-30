import { div, style, State, Fmt, img, a, h2, span } from "./html.js";
import { getCategory } from "./categories.js";
import { Player } from "./player.js";

export class Main extends State {
  constructor(main) {
    let { $, listen, f } = super();

    $("search", document.querySelector("#search"));
    $("submit", document.querySelector("#submit"));
    $("pages", []);
    $("stories", document.querySelector("#stories"));
    $("suggestedCont", document.querySelector("#suggestedCont"));
    $("cats", document.querySelector("#cats"));
    $("suggested", []);
    (() => {
      let { pages } = $();

      for (const key in localStorage) {
        if (
          Object.hasOwnProperty.call(localStorage, key) &&
          key != "settings"
        ) {
          const element = localStorage[key];
          pages.push(JSON.parse(element));
        }
      }
      $("pages", pages);
      getSuggested().then((posts) => {
        $("suggested", posts);
      });
    })();

    listen("submit", "click", ({ search, pages }) => {
      fetch(`https://cors.lowsh.workers.dev/?${search.value}`)
        .then((e) => e.text())
        .then((res) => {
          // Create a temporary container element to parse the HTML
          let parser = new DOMParser();
          let doc = parser.parseFromString(res, "text/html");

          // Use Readability on the new document
          let readabilityDoc = new Readability(doc).parse();

          let cont = div``;
          cont.innerHTML = readabilityDoc.content;

          let paragraphs = [...cont.querySelectorAll("p")].map((p) => {
            return p.innerText;
          });
          let page = {
            id: new Date().getTime(),
            title: readabilityDoc.title,
            byline: readabilityDoc.byline,
            excerpt: readabilityDoc.excerpt,
            readingTime: readabilityDoc.length / (236 * 5),
            date: readabilityDoc.publishedTime,
            site: readabilityDoc.siteName,
            catagory: getCategory(readabilityDoc.excerpt),
            image: (
              doc.querySelector("meta[property='og:image']") || { content: "" }
            ).content,
            text: paragraphs.join("\n"),
          };
          localStorage.setItem(page.id, JSON.stringify(page));
          pages.push(page);
          search.value = "";
          $("pages", pages);
        });
    });

    f(({ pages, stories }) => {
      stories.innerHTML = "";
      for (let i = 0; i < pages.length; i++) {
        const page = pages[i];
        let el = Fmt`${div`class="${css.story}"`}
                        ${div`innerText="${
                          page.site
                        }/${page.catagory.toUpperCase()}" class="${css.site}"`}
                        ${a`innerText="${page.title}" class="${css.title}"`}
                        ${div`class="${css.info}"`}
                            ${div`innerText="${page.byline}" class="${css.byline}"`}
                            ${div`innerText="${new Date(page.date)
                              .toString()
                              .split(" ")
                              .slice(0, 4)
                              .join(" ")}" class="${css.date}"`}
                        ${div``}
                            ${div`innerText="Left: ${parseInt(
                              (page.readingTime / page.text.split(" ").length) *
                                (page.text.split(" ").length - page.current)
                            )} min(s)" class="${css.time}"`}
                        ${img`src="${page.image}" class="${css.image}"`}
                    `;
        el.addEventListener("click", () => {
          new Player(page);
        });
        stories.appendChild(el);
      }
    });

    f(({ suggested, suggestedCont, cats }) => {
      suggestedCont.innerHTML = "";
      let used = [];
      for (let i = 0; i < suggested.length; i++) {
        const suggest = suggested[i];
        let el = Fmt`${div`class="${css.story}" style="margin-bottom: 0"`}
                        ${div`innerText="${
                          suggest.site
                        }/${suggest.catagory.toUpperCase()}" class="${
                          css.site
                        }"`}
                        ${a`innerText="${suggest.title}" class="${css.title}"`}
  
                    `;
        el.addEventListener("click", () => {
          new Preview(page);
        });
        suggestedCont.appendChild(el);
        if (used.indexOf(suggest.catagory) == -1) {
          let c = span`innerText="${suggest.catagory}" class="${css.cats}"`;
          cats.appendChild(c);
          used.push(suggest.catagory);
        }
      }
    });
  }
}

let css = style(/* css */ `
.story {
    width: 700px;
    margin: auto;
    margin-bottom: 50px;
    border-radius: 10px;
    padding: 20px;
    max-width: 85%;
    cursor: pointer;
}
.site{
    color: #7a7a7a;
    font-size: 11px;
    margin-bottom: 10px;
}
.title{    
    margin-bottom: 3px;
}
.info{
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    justify-content: space-between;
    color: #969696;
    font-size: 11px;
    font-family: monospace;
}
.time{
    float: right;
    margin-bottom: 12px;
    font-size: 12px;
    font-family: monospace;
}
.image{
    width: 100%;
    border-radius: 10px;
}
.cats {
    padding: 3px;
    border: 2px solid #6ea2ff;
    border-radius: 5px;
    margin: 0px 5px;
    cursor: pointer;
    text-transform: uppercase;
    font-weight: 700;
    color: #2c2c2c;
    font-size: 13px;
}
`);

function getSuggested() {
  return new Promise((resolve, reject) => {
    fetch(
      `https://cors.lowsh.workers.dev/?https://www.reddit.com/r/Longreads/top.json?t=month`
    )
      .then((e) => e.json())
      .then((d) => {
        let posts = d.data.children.map((e) => {
          let p = e.data;
          return {
            title: p.title,
            url: p.url,
            site: p.domain,
            catagory: getCategory(p.title),
          };
        });
        resolve(posts);
      });
  });
}
