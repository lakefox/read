import { style, State, div, a, img, Fmt } from "./html.js";

export class Preview extends State {
  constructor(suggestion) {
    let { $, f } = super();

    let story = div`class="${css.preview}"`;
    $("story", story);
    $("page", {});

    fetch(`https://cors.lowsh.workers.dev/?${suggestion.url}`)
      .then((e) => e.text())
      .then((res) => {
        // Create a temporary container element to parse the HTML
        let parser = new DOMParser();
        let doc = parser.parseFromString(res, "text/html");

        // Use Readability on the new document
        let readabilityDoc = new Readability(doc).parse();

        let cont = div``;
        cont.innerHTML = readabilityDoc.content;

        let imgs = cont.querySelectorAll("img");

        for (let i = 0; i < imgs.length; i++) {
          const element = imgs[i];
          element.setAttribute("width", "");
          element.setAttribute("height", "");
        }

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
          cont: cont,
          url: suggestion.url,
        };
        console.log(page);

        $("page", page);
      });

    f(({ page, story }) => {
      story.innerHTML = "";
      document.body.style.overflow = "hidden";
      let el = Fmt`${div`class="${css.story}"`}
                        ${div`class="${css.row}"`}
                            ${div`innerText="${
                              page.site
                            } / ${page.catagory.toUpperCase()}" class="${
                              css.site
                            }"`}
                            ${div`innerText="X" class="${css.close}"`.on(
                              "click",
                              () => {
                                document.body.removeChild(story);
                                document.body.style.overflow = "";
                              }
                            )}
                        ${a`innerText="${page.title}" class="${css.title}"`}
                        ${div`class="${css.info}"`}
                            ${div`innerText="${page.byline}" class="${css.byline}"`}
                            ${div`innerText="${new Date(page.date)
                              .toString()
                              .split(" ")
                              .slice(0, 4)
                              .join(" ")}" class="${css.date}"`}
                        ${div``}
                            ${div`innerText="${parseInt(
                              (page.readingTime / page.text.split(" ").length) *
                                page.text.split(" ").length
                            )} min(s)" class="${css.time}"`}
                        ${img`src="${page.image}" class="${css.image}"`}
                        ${page.cont}
                        ${div`innerText="Add to Library" class="${css.save}"`.on(
                          "click",
                          () => {
                            document.querySelector("#search").value =
                              suggestion.url;
                            document.querySelector("#submit").click();
                            document.querySelector("#search").value = "";
                            document.body.removeChild(story);
                            document.body.style.overflow = "";
                          }
                        )}
                    `;
      story.appendChild(el);
    });

    document.body.append(story);
  }
}

let css = style(/* css */ `
.preview {
    position: fixed;
    top: calc(50% - 33px);
    left: 50%;
    transform: translate(-50%, -50%);
    width: calc(100% + 1px);
    height: calc(100% - 65px);
    background: #222;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    color: #fff;
    overflow-y: auto;
}
.story {
    width: 700px;
    margin: auto;
    margin-bottom: 50px;
    border-radius: 10px;
    max-width: 85%;
    padding-top: 40px;
    padding-bottom: 40px;
}

.site{
    color: #d2d2d2;
    font-size: 14px;
    margin-bottom: 10px;
}
.title{    
    margin-bottom: 3px;
}
.row{
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    justify-content: space-between;
    align-items: center;
}
.info{
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    justify-content: space-between;
    color: #e1e1e1;
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
    margin-top: 20px;
}
.save {
    width: 100%;
    height: 50px;
    border-radius: 10px;
    color: #fff;
    background: #88a8ff;
    text-align: center;
    line-height: 50px;
    font-weight: 900;
    cursor: pointer;
    position: sticky;
    bottom: 6px;
}
.close {
    color: #f88;
    font-weight: 900;
    cursor: pointer;
}
`);
