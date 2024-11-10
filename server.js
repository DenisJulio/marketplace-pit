var bs = require("browser-sync").create();

bs.init({
    proxy: {
        target: "http://localhost:3030",
    },
    port: 7070,
})
