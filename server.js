var bs = require("browser-sync").create();

bs.init({
    proxy: {
        target: "http://localhost:7000",
    },
    port: 7070,
})
