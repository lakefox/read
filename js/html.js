export let a = tempToDOM("a");
export let abbr = tempToDOM("abbr");
export let address = tempToDOM("address");
export let area = tempToDOM("area");
export let article = tempToDOM("article");
export let aside = tempToDOM("aside");
export let audio = tempToDOM("audio");
export let b = tempToDOM("b");
export let base = tempToDOM("base");
export let bdi = tempToDOM("bdi");
export let bdo = tempToDOM("bdo");
export let blockquote = tempToDOM("blockquote");
export let body = tempToDOM("body");
export let br = tempToDOM("br");
export let button = tempToDOM("button");
export let canvas = tempToDOM("canvas");
export let caption = tempToDOM("caption");
export let cite = tempToDOM("cite");
export let code = tempToDOM("code");
export let col = tempToDOM("col");
export let colgroup = tempToDOM("colgroup");
export let data = tempToDOM("data");
export let datalist = tempToDOM("datalist");
export let dd = tempToDOM("dd");
export let del = tempToDOM("del");
export let details = tempToDOM("details");
export let dfn = tempToDOM("dfn");
export let dialog = tempToDOM("dialog");
export let div = tempToDOM("div");
export let dl = tempToDOM("dl");
export let dt = tempToDOM("dt");
export let em = tempToDOM("em");
export let embed = tempToDOM("embed");
export let fieldset = tempToDOM("fieldset");
export let figcaption = tempToDOM("figcaption");
export let figure = tempToDOM("figure");
export let footer = tempToDOM("footer");
export let form = tempToDOM("form");
export let h1 = tempToDOM("h1");
export let h2 = tempToDOM("h2");
export let h3 = tempToDOM("h3");
export let h4 = tempToDOM("h4");
export let h5 = tempToDOM("h5");
export let h6 = tempToDOM("h6");
export let head = tempToDOM("head");
export let header = tempToDOM("header");
export let hr = tempToDOM("hr");
export let html = tempToDOM("html");
export let i = tempToDOM("i");
export let iframe = tempToDOM("iframe");
export let img = tempToDOM("img");
export let input = tempToDOM("input");
export let ins = tempToDOM("ins");
export let kbd = tempToDOM("kbd");
export let label = tempToDOM("label");
export let legend = tempToDOM("legend");
export let li = tempToDOM("li");
export let link = tempToDOM("link");
export let main = tempToDOM("main");
export let map = tempToDOM("map");
export let mark = tempToDOM("mark");
export let meta = tempToDOM("meta");
export let meter = tempToDOM("meter");
export let nav = tempToDOM("nav");
export let noscript = tempToDOM("noscript");
export let object = tempToDOM("object");
export let ol = tempToDOM("ol");
export let optgroup = tempToDOM("optgroup");
export let option = tempToDOM("option");
export let output = tempToDOM("output");
export let p = tempToDOM("p");
export let param = tempToDOM("param");
export let picture = tempToDOM("picture");
export let pre = tempToDOM("pre");
export let progress = tempToDOM("progress");
export let q = tempToDOM("q");
export let rp = tempToDOM("rp");
export let rt = tempToDOM("rt");
export let ruby = tempToDOM("ruby");
export let s = tempToDOM("s");
export let samp = tempToDOM("samp");
export let script = tempToDOM("script");
export let section = tempToDOM("section");
export let select = tempToDOM("select");
export let slot = tempToDOM("slot");
export let small = tempToDOM("small");
export let source = tempToDOM("source");
export let span = tempToDOM("span");
export let strong = tempToDOM("strong");
export let sub = tempToDOM("sub");
export let summary = tempToDOM("summary");
export let sup = tempToDOM("sup");
export let table = tempToDOM("table");
export let tbody = tempToDOM("tbody");
export let td = tempToDOM("td");
export let template = tempToDOM("template");
export let textarea = tempToDOM("textarea");
export let tfoot = tempToDOM("tfoot");
export let th = tempToDOM("th");
export let thead = tempToDOM("thead");
export let time = tempToDOM("time");
export let title = tempToDOM("title");
export let tr = tempToDOM("tr");
export let track = tempToDOM("track");
export let u = tempToDOM("u");
export let ul = tempToDOM("ul");
export let video = tempToDOM("video");
export let wbr = tempToDOM("wbr");

export let style = function () {
  let css = templateToString(...arguments);
  let newCSS = renameCSSClasses(css);
  let s = document.createElement("style");
  s.innerHTML = newCSS.modifiedCSS;
  document.head.appendChild(s);
  return newCSS.classMap;
};

export let Fmt = function (strings, ...values) {
  const elements = strings.map((e, i) => {
    return { el: values[i], depth: e.length };
  });

  const stack = [elements[0]];

  for (let i = 1; i < elements.length; i++) {
    const { el, depth } = elements[i];

    if (el == undefined) {
      break;
    }

    while (stack.length > 1 && depth <= stack[stack.length - 1].depth) {
      stack.pop();
    }

    stack[stack.length - 1].el.appendChild(el);
    stack.push({ depth, el });
  }

  return values[0];
};

function tempToDOM(type) {
  return function () {
    let el = document.createElement(type);
    let props = parseHTMLProperties(...arguments);
    for (let prop of Object.keys(props)) {
      if (el[prop] !== undefined) {
        el[prop] = props[prop];
      } else {
        el.setAttribute(prop, props[prop]);
      }
    }
    el.on = (evt, handler) => {
      el.addEventListener(evt, handler);
      return el;
    };
    el.bind = (state, name) => {
      if (typeof el.value != "undefined") {
        el.addEventListener("input", () => {
          state.val(name, el.value);
        });
        state.f((d) => {
          el.value = d[name];
        });
      } else {
        state.f((d) => {
          el.innerHTML = d[name];
        });
      }

      return el;
    };
    el.add = el.appendChild;
    el.clear = () => {
      el.innerHTML = "";
    };
    return el;
  };
}
function parseHTMLProperties(template, ...values) {
  const properties = {};
  const templateString = template.reduce((acc, part, i) => {
    const value = values[i];
    if (typeof value === "string") {
      acc += part + value;
    } else {
      // If the value is not a string, treat it as a property name
      const propName = value;
      if (propName !== undefined) {
        acc += part + `${propName}`;
      } else {
        acc += part;
      }
    }
    return acc;
  }, "");
  // Regular expression to match attribute="value" or attribute='value' or attribute
  const attributeRegex = /([\w-]+)(?:=(?:"([^"]*)"|'([^']*)'))?/g;
  templateString.replace(
    attributeRegex,
    (match, name, doubleQuotedValue, singleQuotedValue) => {
      const value = doubleQuotedValue || singleQuotedValue || true;

      properties[name] = value;
    }
  );

  return properties;
}

function templateToString(strings, ...values) {
  let result = "";
  for (let i = 0; i < strings.length; i++) {
    result += strings[i];
    if (i < values.length) {
      result += values[i];
    }
  }
  return result;
}

function generateRandomHash(length) {
  let result = "";
  const characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const charactersLength = characters.length;
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}

function renameCSSClasses(css) {
  const classMap = {};
  const cssWithModifiedClasses = css.replace(
    /(\.)([\w-]+(?![^{}]*\}))/g,
    (match, selectorType, selectorName) => {
      if (selectorType === "." || selectorType === "#") {
        if (!classMap[selectorName]) {
          classMap[selectorName] = `${selectorName}-${generateRandomHash(8)}`;
        }
        return `${selectorType}${classMap[selectorName]}`;
      } else {
        return match; // Preserve properties like background
      }
    }
  );

  return {
    modifiedCSS: cssWithModifiedClasses,
    classMap,
  };
}

export class State {
  constructor() {
    let auth = false;
    let context = {};
    let dependencies = {};
    let renders = {};
    this.f = async (render = () => {}) => {
      let name = `f${Object.keys(renders).length}`;
      renders[name] = render;
      try {
        await run(render, name, context);
      } catch (error) {
        // console.error(error);
      }
    };
    this.listen = (name, ...events) => {
      let handler = events.at(-1);
      for (let i = 0; i < events.length - 1; i++) {
        const event = events[i];
        context[name].addEventListener(event, async () => {
          let ctx = context;
          ctx.parent = context[name];
          await run(handler, name, ctx);
        });
      }
    };
    this.$ = (name, value) => {
      if (value === undefined && name === undefined) {
        return context;
      } else if (value === undefined && typeof name == "string") {
        return context[name];
      } else if (typeof name == "object") {
        for (const key in name) {
          context[key] = name[key];
          if (dependencies[key]) {
            for (let i = 0; i < dependencies[key].length; i++) {
              let render = renders[dependencies[key][i]];
              if (render) {
                render(context);
              }
            }
          }
        }
        if (auth) {
          console.log(context);
          auth.ch.emit(auth.type, serialize(context));
        }
      } else {
        context[name] = value;
        if (dependencies[name]) {
          for (let i = 0; i < dependencies[name].length; i++) {
            let render = renders[dependencies[name][i]];
            if (render) {
              render(context);
            }
          }
        }
        if (auth) {
          console.log(context);
          auth.ch.emit(auth.type, serialize(context));
        }
      }
    };
    this.mount = (type, ch, restore = false) => {
      auth = { type, ch };
      if (restore) {
        ch.emit(type);
        ch.on(type, async (data) => {
          console.log(type, data);
          if (data != null) {
            console.log("not null");
            context = unSerialize(data);
            // Load the data and render the current state
            for (const key in context) {
              for (const r in renders) {
                await run(renders[r], key, context);
              }
            }
          }
        });
      }
    };
    let state = this;
    this.Global = class {
      constructor() {
        let events = {};
        let catchAll = () => {};
        this.call = (event, ...args) => {
          if (events[event]) {
            for (let i = 0; i < events[event].length; i++) {
              events[event][i](...args);
            }
          } else {
            catchAll(event, ...args);
          }
        };
        this.$ = state.$;
        this.listen = state.listen;
        this.f = state.f;
        this.mount = state.mount;
        this.on = (event, cb) => {
          if (!events[event]) {
            events[event] = [];
          }
          events[event].push(cb);
        };
        this.catch = (cA) => {
          catchAll = cA;
        };
      }
    };
    async function run(func, name, context) {
      let require = {};
      for (const key in context) {
        if (Object.hasOwnProperty.call(context, key)) {
          Object.defineProperty(require, key, {
            get() {
              if (dependencies[key] == undefined) {
                dependencies[key] = [];
              }
              if (dependencies[key].indexOf(name) == -1) {
                dependencies[key].push(name);
              }
              return context[key];
            },
          });
        }
      }
      let data = func(require);
      if (typeof data?.then === "function") {
        data = await data;
      }
      Object.assign(context, data);
      for (const key in dependencies) {
        if (Object.hasOwnProperty.call(dependencies, key)) {
          for (let i = 0; i < dependencies[key].length; i++) {
            if (renders[dependencies[key][i]] && dependencies[key][i] != name) {
              renders[dependencies[key][i]](context);
            }
          }
        }
      }
    }
  }
}
