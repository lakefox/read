import { div, style, State, Fmt, img, a } from "./html.js";
import { getCategory } from "./categories.js";
import { Dialog } from "./dialog.js";

export class Main extends State {
  constructor(main) {
    let { val, listen, f } = super();

    val("search", document.querySelector("#search"));
    val("submit", document.querySelector("#submit"));
    val("pages", []);
    val("stories", document.querySelector("#stories"));
    (() => {
      let { pages } = val();

      for (const key in localStorage) {
        if (Object.hasOwnProperty.call(localStorage, key)) {
          const element = localStorage[key];
          pages.push(JSON.parse(element));
        }
      }
      val("pages", pages);
    })();

    listen("submit", "click", ({ search, pages }) => {
      console.log(search.value);
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
            readingTime: readabilityDoc.length / 120,
            date: readabilityDoc.publishedTime,
            site: readabilityDoc.siteName,
            catagory: getCategory(readabilityDoc.excerpt),
            image: (
              doc.querySelector("meta[property='og:image']") || { content: "" }
            ).content,
            paragraphs,
            current: 0,
          };
          localStorage.setItem(page.id, JSON.stringify(page));
          pages.push(page);
          search.value = "";
          val("pages", pages);
        });
    });

    f(({ pages, stories }) => {
      console.log(pages);
      stories.innerHTML = "";
      for (let i = 0; i < pages.length; i++) {
        const page = pages[i];
        console.log();
        let el = Fmt`${div`class="${css.story}"`}
                        ${div`innerText="${
                          page.site
                        }/${page.catagory.toUpperCase()}" class="${css.site}"`}
                        ${a`innerText="${page.title}" class="${css.title}"`}
                        ${div`class="${css.info}"`}
                            ${div`innerText="${page.byline}" class="${css.byline}"`}
                            ${div`innerText="${page.date}" class="${css.date}"`}
                        ${div``}
                            ${div`innerText="Left: ${parseInt(
                              (page.readingTime / page.paragraphs.length) *
                                (page.paragraphs.length - page.current)
                            )} min(s)" class="${css.time}"`}
                        ${img`src="${page.image}" class="${css.image}"`}
                    `;
        el.addEventListener("click", () => {
          console.log(page);
          new Dialog(page);
        });
        console.log(el);
        stories.appendChild(el);
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

`);