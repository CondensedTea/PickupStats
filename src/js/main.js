function createRatingList(items, type) {
    const rating = document.createElement('tbody');
    rating.id = "tbody"
    document.getElementById('render').appendChild(rating);

    items.forEach((item, index) => {
        let tr = rating.insertRow(-1);

        let th = document.createElement('th');
        th.setAttribute('scope', 'row')
        let indexText = document.createTextNode(index + 1);
        th.appendChild(indexText)
        tr.appendChild(th)

        let cellName = tr.insertCell(-1);
        // let name = document.createTextNode(item.player_name);
        let a = document.createElement('a')
        a.text = (item?.player_name === "") ? item.steamid64 : item.player_name;
        a.href = "https://steamcommunity.com/profiles/" + item.steamid64
        cellName.appendChild(a)

        switch (type) {
            case 'kdr':
                let cellKDR = tr.insertCell(-1);
                let kdr = document.createTextNode(item.kdr)
                cellKDR.appendChild(kdr)
                break
            case 'dpm':
                let cellDPM = tr.insertCell(-1);
                let dpm = document.createTextNode(item.dpm)
                cellDPM.appendChild(dpm)
                break
            case 'hpm':
                let cellHPM = tr.insertCell(-1);
                let hpm = document.createTextNode(item.hpm)
                cellHPM.appendChild(hpm)
                break
        }
        let cellGamesCount = tr.insertCell(-1);
        let games = document.createTextNode(item.games);
        cellGamesCount.appendChild(games)
    });
}

async function updateRatingList(url, type) {
    document.getElementById("tbody").remove()
    let playerClassElem = document.getElementById("playerClass")

    let playerClass
    if (type === "hpm") {
        playerClass = ""
    } else {
        let playerClass = playerClassElem.options[playerClassElem.selectedIndex].value

    }
    let minGamesElem = document.getElementById("minGames")
    let mingames = minGamesElem.value

    let params
    if (playerClass === "Any" || playerClass === "") {
        params = new URLSearchParams({'mingames': mingames})
    } else {
        params = new URLSearchParams({'class': playerClass.toLowerCase(), 'mingames': mingames})
    }
    let items = await getDataFromAPI(`${url}?${params.toString()}`)
    createRatingList(items['stats'], type)
}

async function updateGameCount() {
    let data = await getDataFromAPI('/api/gamesCount')
    let elem = document.getElementById("gamesCounter")

    elem.innerText += " " + data["count"]
}

async function getDataFromAPI(url) {
    let resp = await fetch(url);
    return await resp.json();
}
