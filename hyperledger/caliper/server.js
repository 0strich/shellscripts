const express = require("express");
const path = require("path");
const app = express();

// 정적 파일 제공
app.use(express.static(path.join(__dirname)));

// 8080 포트에서 서버 시작
app.listen(3000, () => {
  console.log("Server is running at http://localhost:3000");
});
