import {
  div,
  style,
  State,
  Fmt,
  img,
  input,
  button,
  progress,
} from "./html.js";
import { textToSpeech } from "./tts.js";
window.speechSynthesis.cancel();

export class Dialog extends State {
  constructor(page) {
    let { val, listen, f } = super();

    let timeOut;

    val("body", document.body);
    val("play", true);
    val("current", page.current);
    val("sleep", 0);

    let text = div`class="${css.text}" innerText="${
      page.paragraphs[page.current]
    }"`;
    let controls = button`class="${css.controls}"`.on("click", () => {
      val("play", !val("play"));
    });
    let back = button`class="${css.dir}" innerHTML="&#9664;&#9664;"`.on(
      "click",
      () => {
        window.speechSynthesis.cancel();

        val("current", Math.max(val("current") - 1, 0));
      }
    );
    let forward = button`class="${css.dir}" innerHTML="&#9654;&#9654;"`.on(
      "click",
      () => {
        window.speechSynthesis.cancel();

        val(
          "current",
          Math.min(val("current") + 1, page.paragraphs.length - 1)
        );
      }
    );

    let p = progress`value="${page.current + 5}" max="${
      page.paragraphs.length
    }" class="${css.progress}"`;

    let d = Fmt`${div`class="${css.dialog}"`}
                    ${img`src="${page.image}" class="${css.bg}"`}
                    ${div`class="${css.row}"`}
                        ${div`innerText="${page.title}" class="${css.title}"`}
                        ${div`innerText="X" class="${css.exit}"`.on(
                          "click",
                          () => {
                            window.speechSynthesis.cancel();
                            document.querySelector("." + css.dialog).remove();
                          }
                        )}
                    ${div`class="${css.cont}"`}
                        ${text}
                    ${div`class="${css.bottom}"`}
                        ${div``}
                            ${div`class="${css.sleep}" innerHTML="&#10088;"`.on(
                              "click",
                              () => {
                                prompt2("Time Until Sleep", "text", 10).then(
                                  (time) => {
                                    val("sleep", parseInt(time));
                                  }
                                );
                              }
                            )}
                            ${back}
                            ${controls}
                            ${forward}
                            ${div`class="${css.settings}" innerHTML="&#9881;"`}
                        ${p}

                    `;
    // timeout 5/15/30/1 hour
    document.body.appendChild(d);

    play(page, val);

    f((e) => {
      if (e.play) {
        controls.innerHTML = "&#10074;&#10074;";
        play(page, val);
        if (e.sleep > 0) {
          timeOut = setTimeout(() => {
            window.speechSynthesis.cancel();
            val("play", false);
          }, e.sleep * 600);
        }
      } else {
        controls.innerHTML = "&#9654;";
        window.speechSynthesis.cancel();
      }
    });

    f(({ current }) => {
      p.value = current;
      page.current = current;
      text.innerText = page.paragraphs[page.current];
      localStorage.setItem(page.id, JSON.stringify(page));
      if (val("play")) {
        console.log(val("current"));
        play(page, val);
      }
    });

    f(({ sleep }) => {
      if (sleep) {
        clearTimeout(timeOut);
        timeOut = setTimeout(() => {
          window.speechSynthesis.cancel();
          val("play", false);
        }, sleep * 600);
      }
    });

    function play(page, val) {
      textToSpeech(
        page.paragraphs[page.current],
        {
          lang: "en-US",
          voice: "Nicky",
          pitch: 1.2,
          rate: 1,
          volume: 0.8,
        },
        () => {
          page.current++;
          val("current", page.current);
        }
      );
    }

    // listen("submit", "click", ({ search, pages }) => {});
  }
}

let css = style(/* css */ `
.dialog {
    width: 100vw;
    height: 100vh;
    position: absolute;
    top: 0;
    left: 0;
    background: #fff;
    overflow: hidden;
}
.bg {
    width: 100vw;
    height: 100vh;
    filter: blur(35px);
}


.title {
    position: fixed;
    top: 25px;
    left: 25px;
    color: #fff;
    font-family: sans-serif;
    font-weight: 700;
    font-size: 18px;
    max-width: 650px;
    width: 70%;
}

.exit {
    position: fixed;
    top: 10px;
    right: 20px;
    color: #2a2a2a;
    font-weight: 900;
    font-size: 27px;
    cursor: pointer;
}
.row{}
.bg{}
.cont {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    color: #000;
    font-size: 23px;
    width: 700px;
    max-width: 86%;
    padding: 15px;
    border-radius: 10px;
    font-weight: 900;
}
.text {
    text-align: center;
    max-height: 500px;
    overflow-y: auto;
}
.bottom {
    position: fixed;
    bottom: 0;
    width: 100%;
    height: 120px;
    display: flex;
    flex-direction: column;
    align-items: center;
}
.progress {
    width: 90%;
}
.controls {
    font-size: 40px;
    background: transparent;
    border: none;
    color: #fff;
}
.dir {
    font-size: 33px;
    background: transparent;
    border: none;
    color: #fff;
}
.settings {
    position: absolute;
    right: 6%;
    bottom: 77px;
    font-size: 37px;
    color: #000000ab;
}
.sleep {
    position: absolute;
    left: 6%;
    bottom: 77px;
    font-size: 37px;
    color: #000000ab;
    cursor: pointer;
}
.prompt {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: #212121;
    padding: 20px;
    border-radius: 6px;
    width: 85%;
    height: 180px;
    color: #fff;
    display: flex;
    flex-direction: column;
    align-items: center;
}
.ptext {
    font-size: 22px;
}
.pinput {
    width: 90%;
    height: 20px;
    border: none;
    border-radius: 2px;
    margin: 15px;
}
.prow {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
    position: absolute;
    bottom: 45px;
}
.pbtn2{
    width: 140px;
    background: #606060;
    height: 35px;
    border-radius: 3px;
    cursor: pointer;
    text-align: center;
    line-height: 35px;
    margin-right: 10px;
}
.pbtn{
    width: 140px;
    background: #6685ff;
    height: 35px;
    border-radius: 3px;
    cursor: pointer;
    text-align: center;
    line-height: 35px;
    margin-left: 10px;
}

`);

function prompt2(text, type = "text", value = "") {
  return new Promise((resolve, reject) => {
    let d = Fmt`${div`class="${css.prompt}"`}
                    ${div`class="${css.ptext}" innerText="${text}"`}
                    ${input`class="${css.pinput}" type="${type}" value="${value}"`}
                    ${div`class="${css.prow}"`}
                        ${div`class="${css.pbtn2}" innerText="Back"`.on(
                          "click",
                          () => {
                            d.remove();

                            reject();
                          }
                        )}
                        ${div`class="${css.pbtn}" innerText="Submit"`.on(
                          "click",
                          (e) => {
                            d.remove();
                            resolve(
                              e.target.parentNode.parentNode.querySelector(
                                "input"
                              ).value
                            );
                          }
                        )}
                `;
    document.body.appendChild(d);
  });
}