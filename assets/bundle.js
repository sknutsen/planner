export function show(elemId) {
  let elem = document.getElementById(elemId);

  elem.classList.remove("hidden");
}

export function hide(elemId) {
  let elem = document.getElementById(elemId);

  elem.classList.add("hidden");
}
