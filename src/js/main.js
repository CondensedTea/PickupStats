function createRatingList(items, type) {
    const rating = document.createElement('ol');
    document.getElementById('render').appendChild(rating);

    items.forEach((item) => {
        let li = document.createElement('li');
        switch (type) {
            case 'kdr':
                li.innerHTML = `<i> ${item.steamid64} | kdr: ${item.kdr}, games: ${item.games} </i>`;
                break
            case 'dpm':
                li.innerHTML = `<i> ${item.steamid64} | dpm: ${item.dpm}, games: ${item.games} </i>`;
                break
            case 'hpm':
                li.innerHTML = `<i> ${item.steamid64} | hpm: ${item.hpm}, games: ${item.games} </i>`;
                break
        }
        rating.appendChild(li);
    });
}

async function getRatingsFromAPI(url) {
    let resp = await fetch(url);
    return await resp.json();
}

