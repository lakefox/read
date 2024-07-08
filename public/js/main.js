import { div, style, State, Fmt, img, a, h2, span } from "./html.js";
import { Player } from "./player.js";
import { Preview } from "./preview.js";

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
    $("filter", []);
    $("library", "mason");

    (() => {
      getSuggested().then((posts) => {
        $("suggested", posts);
      });
    })();

    f(({ library }) => {
      fetch(`/library/${library}`)
        .then((e) => e.json())
        .then((res) => {
          console.log(res);
          $("pages", res || []);
        });
    });

    listen("submit", "click", ({ search, library }) => {
      let url = search.value;
      fetch(`/library/${library}`, {
        method: "POST",
        body: JSON.stringify({ url }),
      })
        .then((e) => e.json())
        .then((res) => {
          search.value = "";
          $("pages", res);
        });
    });

    f(({ pages, stories }) => {
      stories.innerHTML = "";
      for (let i = 0; i < pages.length; i++) {
        let page = pages[pages.length - 1 - i];
        console.log(page);
        let el = Fmt`${div`class="${css.story}"`}
                        ${div`innerText="${page.site.toUpperCase()} / ${page.category.toUpperCase()}" class="${
                          css.site
                        }"`}
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
                                page.text.split(" ").length
                            )} min(s)" class="${css.time}"`}
                        ${img`src="${page.image}" class="${css.image}"`}
                    `;
        el.addEventListener("click", () => {
          new Player(page);
        });
        stories.appendChild(el);
      }
    });

    f(({ suggested, suggestedCont, cats, filter }) => {
      suggestedCont.innerHTML = "";
      cats.innerHTML = "";
      let used = [];
      for (let i = 0; i < suggested.length; i++) {
        let suggest = suggested[i];
        if (filter.length == 0 || filter.indexOf(suggest.category) != -1) {
          let el = Fmt`${div`class="${css.story}" style="margin-bottom: 0"`}
                        ${div`innerText="${suggest.site
                          .split(".")[0]
                          .toUpperCase()} / ${(
                          suggest.category || ""
                        ).toUpperCase()}" class="${css.site}"`}
                        ${a`innerText="${suggest.title}" class="${css.title}"`}
  
                    `;
          el.addEventListener("click", () => {
            new Preview(suggest);
          });
          suggestedCont.appendChild(el);
          if (
            used.indexOf(suggest.category) == -1 &&
            suggest.category != null
          ) {
            let c =
              span`innerText="${suggest.category}" class="${css.cats}"`.on(
                "click",
                () => {
                  let { filter } = $();
                  filter.push(suggest.category);
                  $("filter", filter);
                }
              );
            cats.appendChild(c);
            used.push(suggest.category);
          }
        }
      }
      if (filter.length > 0) {
        let c = span`innerText="X" class="${css.clearcats}"`.on("click", () => {
          let { filter } = $();
          filter = [];
          $("filter", filter);
        });
        cats.appendChild(c);
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
.clearcats {
    padding: 3px;
    padding-left: 5px;
    padding-right: 5px;
    border: 2px solid #989898;
    border-radius: 5px;
    margin: 0px 5px;
    cursor: pointer;
    text-transform: uppercase;
    font-weight: 700;
    color: #989898;
    font-size: 13px; 
}
`);

function getSuggested() {
  let keys = Object.keys(localStorage);
  let urls = [];
  for (let i = 0; i < keys.length; i++) {
    urls.push(JSON.parse(localStorage[keys[i]]).url);
  }
  return new Promise((resolve, reject) => {
    fetch(`https://api.szn.io/feed`, {
      method: "POST",
      body: JSON.stringify({
        data: urls,
      }),
    })
      .then((e) => e.json())
      .then((d) => {
        let posts = d.data.map((p) => {
          return {
            title: p.title,
            url: p.url,
            site: p.site,
            category: p.category,
          };
        });

        resolve(posts);
      });
  });
}
