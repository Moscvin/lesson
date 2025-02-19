import fs from "fs";

fs.appendFile("my-file.txt", "Fail create", (err) => {
  if (err) throw err;
  console.log("Saved!");
});

setTimeout(() => {
  console.log("Hello");
}, 10000);
