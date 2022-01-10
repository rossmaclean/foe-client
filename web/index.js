const express = require('express');
const bodyParser = require('body-parser');

const app = express();

app.use(bodyParser.json());

app.post('/report', (req, res) => {
    console.log(req.body);
    res.sendStatus(200);
});

app.listen(8080, () => console.log(`Started server at http://localhost:8080!`));