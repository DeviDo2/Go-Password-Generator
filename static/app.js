const form = document.querySelector("#password-form");
const generateButton = document.querySelector("#generate-now");
const copyAllButton = document.querySelector("#copy-all");
const passwordList = document.querySelector("#password-list");
const errorMessage = document.querySelector("#error-message");
const poolSize = document.querySelector("#pool-size");
const entropy = document.querySelector("#entropy");

const fields = {
  length: document.querySelector("#length"),
  count: document.querySelector("#count"),
  lowercase: document.querySelector("#lowercase"),
  uppercase: document.querySelector("#uppercase"),
  numbers: document.querySelector("#numbers"),
  symbols: document.querySelector("#symbols"),
};

function requestPayload() {
  return {
    length: Number(fields.length.value),
    count: Number(fields.count.value),
    lowercase: fields.lowercase.checked,
    uppercase: fields.uppercase.checked,
    numbers: fields.numbers.checked,
    symbols: fields.symbols.checked,
  };
}

async function generatePasswords() {
  errorMessage.textContent = "";

  const response = await fetch("/api/generate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(requestPayload()),
  });

  const data = await response.json();
  if (!response.ok) {
    passwordList.replaceChildren();
    errorMessage.textContent = data.error || "Unable to generate passwords.";
    poolSize.textContent = "0 chars";
    entropy.textContent = "0 bits";
    return;
  }

  poolSize.textContent = `${data.characterPool} chars`;
  entropy.textContent = `${data.entropyBits} bits`;
  renderPasswords(data.passwords);
}

function renderPasswords(passwords) {
  const items = passwords.map((password) => {
    const item = document.createElement("li");
    item.className = "password-item";

    const value = document.createElement("code");
    value.textContent = password.value;

    const meta = document.createElement("span");
    meta.textContent = `${password.entropyBits} bits`;

    const button = document.createElement("button");
    button.type = "button";
    button.textContent = "Copy";
    button.addEventListener("click", async () => {
      await navigator.clipboard.writeText(password.value);
      button.textContent = "Copied";
      setTimeout(() => {
        button.textContent = "Copy";
      }, 1200);
    });

    item.append(value, meta, button);
    return item;
  });

  passwordList.replaceChildren(...items);
}

async function copyAll() {
  const values = [...passwordList.querySelectorAll("code")].map((node) => node.textContent);
  if (values.length === 0) {
    return;
  }

  await navigator.clipboard.writeText(values.join("\n"));
  copyAllButton.textContent = "Copied";
  setTimeout(() => {
    copyAllButton.textContent = "Copy all";
  }, 1200);
}

form.addEventListener("input", generatePasswords);
generateButton.addEventListener("click", generatePasswords);
copyAllButton.addEventListener("click", copyAll);
generatePasswords();
