function createRatingList(items, type) {
    const rating = document.createElement('ol');
    document.getElementById('render').appendChild(rating);

    items.forEach((item) => {
        let li = document.createElement('li');
        switch (type) {
            case 'kdr':
                li.innerHTML = `${item.player_name} | kdr: ${item.kdr}, games: ${item.games}`;
                break
            case 'dpm':
                li.innerHTML = `${item.player_name} | dpm: ${item.dpm}, games: ${item.games}`;
                break
            case 'hpm':
                li.innerHTML = `${item.player_name} | hpm: ${item.hpm}, games: ${item.games}`;
                break
        }
        rating.appendChild(li);
    });
}

async function getRatingsFromAPI(url) {
    let resp = await fetch(url);
    return await resp.json();
}

