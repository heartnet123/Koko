const fs = require('fs');
fetch('https://api.jikan.moe/v4/anime/1/full')
  .then(res => res.json())
  .then(data => {
    fs.writeFileSync('jikan_sample.json', JSON.stringify(data, null, 2));
    console.log('Saved to jikan_sample.json');
  })
  .catch(err => console.error(err));
