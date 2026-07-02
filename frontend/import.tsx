import React, { useEffect, useState } from 'react';

export default function App() {
   Load Iconify script dynamically to ensure icons render correctly
  useEffect(() = {
    const script = document.createElement('script');
    script.src = 'httpscode.iconify.designiconify-icon1.0.7iconify-icon.min.js';
    script.async = true;
    document.body.appendChild(script);
    return () = {
      document.body.removeChild(script);
    };
  }, []);

  return (
    div className=min-h-screen bg-[#FCFCFD] text-gray-900 font-sans flex selectionbg-indigo-100 selectiontext-indigo-900
      Sidebar 
      main className=flex-1 flex flex-col min-w-0
        Header 
        div className=flex-1 overflow-y-auto pb-12
          div className=max-w-7xl mx-auto px-8 w-full flex flex-col gap-8 mt-4
            HeroSection 
            RecommendedSection 
            GenreSection 
          div
        div
      main
    div
  );
}

function Sidebar() {
  const menuItems = [
    { icon 'solarhome-smile-linear', label 'Home', active true },
    { icon 'solarcompass-linear', label 'Browse' },
    { icon 'solarfire-linear', label 'Trending' },
    { icon 'solarbookmark-linear', label 'Watchlist' },
  ];

  const userItems = [
    { icon 'solaruser-linear', label 'My Profile' },
    { icon 'solarsettings-linear', label 'Settings' },
  ];

  return (
    aside className=w-64 flex-shrink-0 border-r border-gray-10050 flex flex-col pt-8 pb-6 bg-[#FCFCFD] h-screen sticky top-0
      { Brand Logo }
      div className=px-8 mb-10 flex items-center gap-3 cursor-pointer
        div className=w-1.5 h-6 bg-[#635BFF] rounded-fulldiv
        h1 className=text-2xl font-semibold tracking-tighter text-[#1A1A2E]
          KoKo
        h1
      div

      { Main Navigation }
      nav className=flex-1 px-4 space-y-1
        {menuItems.map((item, idx) = (
          a
            key={idx}
            href=#
            className={`flex items-center gap-4 px-4 py-3 rounded-xl transition-all duration-200 ${
              item.active
                 'bg-[#F3EEFF] text-[#635BFF] font-medium shadow-sm'
                 'text-gray-500 hoverbg-gray-50 hovertext-gray-900'
            }`}
          
            iconify-icon 
              icon={item.icon} 
              width=20 
              style={{ strokeWidth '1.5' }}
            iconify-icon
            span className=text-sm{item.label}span
          a
        ))}

        div className=h-6div { Spacer }

        {userItems.map((item, idx) = (
          a
            key={`user-${idx}`}
            href=#
            className=flex items-center gap-4 px-4 py-3 rounded-xl text-gray-500 hoverbg-gray-50 hovertext-gray-900 transition-all duration-200
          
            iconify-icon 
              icon={item.icon} 
              width=20 
              style={{ strokeWidth '1.5' }}
            iconify-icon
            span className=text-sm{item.label}span
          a
        ))}
      nav

      { Continue Watching Widget }
      div className=px-6 mt-auto
        div className=bg-white rounded-2xl p-4 shadow-[0_2px_12px_-4px_rgba(0,0,0,0.05)] border border-gray-100
          h3 className=text-sm font-medium text-gray-800 mb-3 tracking-tightContinue Watchingh3
          div className=relative rounded-xl overflow-hidden mb-3 aspect-video group cursor-pointer
            img 
              src=httpsimages.unsplash.comphoto-1518173946687-a4c8892bbd9fq=80&w=400&auto=format&fit=crop 
              alt=Episode thumbnail
              className=w-full h-full object-cover transition-transform duration-500 group-hoverscale-105
            
            div className=absolute inset-0 bg-black20 flex items-center justify-center opacity-0 group-hoveropacity-100 transition-opacity duration-300
               div className=w-8 h-8 rounded-full bg-white90 backdrop-blur flex items-center justify-center shadow-lg
                  iconify-icon icon=solarplay-bold width=16 class=text-[#635BFF] ml-0.5iconify-icon
               div
            div
          div
          div className=flex justify-between items-baseline mb-1
            p className=text-sm font-medium text-gray-900 tracking-tightEpisode 4p
            span className=text-xs text-gray-400 font-medium60%span
          div
          p className=text-xs text-gray-400 mb-2Drama • Romancep
          div className=w-full bg-gray-100 h-1.5 rounded-full overflow-hidden
            div className=bg-[#635BFF] w-[60%] h-full rounded-fulldiv
          div
        div
      div
    aside
  );
}

function Header() {
  return (
    header className=h-20 px-8 flex items-center justify-between sticky top-0 bg-[#FCFCFD]80 backdrop-blur-md z-40
      { Search Bar }
      div className=relative w-full max-w-lg group
        div className=absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none
          iconify-icon icon=solarmagnifer-linear width=18 class=text-gray-400 group-focus-withintext-[#635BFF] transition-colorsiconify-icon
        div
        input
          type=text
          placeholder=Search for movies, series, genres...
          className=w-full pl-11 pr-4 py-3 bg-white border border-gray-100 rounded-full text-sm text-gray-900 focusoutline-none focusring-2 focusring-[#635BFF]10 focusborder-[#635BFF]30 transition-all shadow-[0_2px_8px_-4px_rgba(0,0,0,0.02)] placeholdertext-gray-400
        
      div

      { User Actions }
      div className=flex items-center gap-6
        button className=relative text-gray-400 hovertext-gray-600 transition-colors
          iconify-icon icon=solarbell-linear width=22iconify-icon
          span className=absolute top-0.5 right-0.5 w-2 h-2 bg-red-400 border-2 border-[#FCFCFD] rounded-fullspan
        button
        
        div className=flex items-center gap-2 cursor-pointer hoveropacity-80 transition-opacity
          img 
            src=httpsimages.unsplash.comphoto-1535713875002-d1d0cf377fdeq=80&w=100&auto=format&fit=crop 
            alt=User Profile 
            className=w-9 h-9 rounded-full object-cover shadow-sm border border-gray-100
          
          iconify-icon icon=solaralt-arrow-down-linear width=16 class=text-gray-400iconify-icon
        div
      div
    header
  );
}

function HeroSection() {
  return (
    div className=relative w-full h-[380px] rounded-[2rem] overflow-hidden group shadow-[0_8px_30px_rgb(0,0,0,0.04)]
      { Background Image & Gradients }
      img 
        src=httpsimages.unsplash.comphoto-1518173946687-a4c8892bbd9fq=80&w=2000&auto=format&fit=crop 
        alt=Romance in The Sunset 
        className=absolute inset-0 w-full h-full object-cover object-center
      
      { Left to right gradient for text legibility }
      div className=absolute inset-0 bg-gradient-to-r from-[#FDFDFE] via-[#FDFDFE]80 to-transparentdiv
      
      { Content Container }
      div className=relative h-full flex flex-col justify-center px-12 max-w-xl
        h2 className=text-4xl font-semibold tracking-tight text-gray-900 leading-[1.1] mb-4
          Romance inbr The Sunset
        h2
        
        div className=flex items-center gap-4 text-sm mb-4
          span className=text-[#635BFF] font-medium bg-[#635BFF]10 px-2 py-0.5 rounded text-xsNew Episodespan
          span className=text-gray-500 font-mediumEvery Fridayspan
        div
        
        p className=text-sm text-gray-500 leading-relaxed mb-8 max-w-sm
          A heartwarming story of love, fate and second chances.
        p
        
        div className=flex items-center gap-3
          button className=flex items-center justify-center gap-2 bg-[#635BFF] hoverbg-[#5249E5] text-white px-6 py-2.5 rounded-full text-sm font-medium transition-all shadow-lg shadow-[#635BFF]20 activescale-95
            iconify-icon icon=solarplay-bold width=16iconify-icon
            Play Now
          button
          button className=flex items-center justify-center gap-2 bg-white hoverbg-gray-50 text-gray-700 px-6 py-2.5 rounded-full text-sm font-medium transition-all shadow-sm border border-gray-100 activescale-95
            iconify-icon icon=solaradd-circle-linear width=18iconify-icon
            Watchlist
          button
        div
      div

      { Slider Controls }
      div className=absolute right-8 top-12 -translate-y-12
         button className=w-10 h-10 bg-white80 backdrop-blur-md rounded-full flex items-center justify-center text-gray-600 hoverbg-white transition-all shadow-sm opacity-0 group-hoveropacity-100 translate-x-4 group-hovertranslate-x-0
            iconify-icon icon=solaralt-arrow-right-linear width=20iconify-icon
         button
      div

      { Pagination Dots }
      div className=absolute bottom-6 left-12 -translate-x-12 flex items-center gap-2
        div className=w-6 h-1.5 bg-[#635BFF] rounded-fulldiv
        div className=w-1.5 h-1.5 bg-white60 rounded-fulldiv
        div className=w-1.5 h-1.5 bg-white60 rounded-fulldiv
        div className=w-1.5 h-1.5 bg-white60 rounded-fulldiv
        div className=w-1.5 h-1.5 bg-white60 rounded-fulldiv
      div
    div
  );
}

function RecommendedSection() {
  const [recommendations, setRecommendations] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() = {
    const fetchRecommendations = async () = {
      try {
        const res = await fetch('httpsapi.jikan.moev4recommendationsanime');
        const data = await res.json();
        
        if (data && data.data) {
           The API returns pairs of recommended anime. 
           We'll extract unique anime entries up to a limit of 5.
          const uniqueAnimes = [];
          const seenIds = new Set();
          
          for (const rec of data.data) {
             Extract the first anime from the recommendation pair
            const anime = rec.entry[0];
            if (anime && !seenIds.has(anime.mal_id)) {
              seenIds.add(anime.mal_id);
              uniqueAnimes.push(anime);
            }
            if (uniqueAnimes.length = 5) break;
          }
          
          setRecommendations(uniqueAnimes);
        }
      } catch (error) {
        console.error(Failed to fetch recommendations, error);
      } finally {
        setLoading(false);
      }
    };
    fetchRecommendations();
  }, []);

  return (
    section
      div className=flex items-center justify-between mb-5
        h3 className=text-lg font-medium tracking-tight text-gray-900Recommendedh3
        button className=text-sm font-medium text-gray-500 hovertext-gray-900 flex items-center gap-1 transition-colors group
          See All 
          iconify-icon icon=solaralt-arrow-right-linear width=16 class=group-hovertranslate-x-0.5 transition-transformiconify-icon
        button
      div
      
      {loading  (
        div className=flex items-center justify-center h-32 text-sm text-gray-400 bg-gray-50 rounded-2xl border border-gray-100 border-dashed
          Loading recommendations...
        div
      )  (
        div className=grid grid-cols-5 gap-5
          {recommendations.map((item) = (
            div key={item.mal_id} className=flex flex-col group cursor-pointer
              div className=relative rounded-2xl overflow-hidden aspect-[43] mb-3 shadow-[0_4px_12px_rgb(0,0,0,0.03)] border border-gray-10050 bg-gray-100
                img 
                  src={item.images.jpg.large_image_url  item.images.jpg.image_url} 
                  alt={item.title}
                  className=w-full h-full object-cover transition-transform duration-700 group-hoverscale-105
                
                { Play Button Overlay }
                div className=absolute bottom-3 left-3 w-8 h-8 rounded-full bg-white90 backdrop-blur shadow-sm flex items-center justify-center text-gray-800 opacity-90 group-hoveropacity-100 group-hoverscale-110 transition-all
                  iconify-icon icon=solarplay-bold width=14 class=ml-0.5iconify-icon
                div
              div
              h4 className=text-sm font-medium text-gray-900 tracking-tight truncate title={item.title}{item.title}h4
              p className=text-xs text-gray-500 mt-0.5Highly Recommendedp
            div
          ))}
        div
      )}
    section
  );
}

const getGenreStyle = (index) = {
  const styles = [
    { icon 'solarbomb-linear', color 'bg-rose-50 text-rose-600' },
    { icon 'solarshield-linear', color 'bg-emerald-50 text-emerald-600' },
    { icon 'solaremoji-funny-circle-linear', color 'bg-amber-50 text-amber-600' },
    { icon 'solarheart-linear', color 'bg-pink-50 text-pink-600' },
    { icon 'solarrocket-linear', color 'bg-sky-50 text-sky-600' },
    { icon 'solarghost-linear', color 'bg-purple-50 text-purple-600' },
    { icon 'solarmagic-stick-3-linear', color 'bg-indigo-50 text-indigo-600' },
    { icon 'solarusers-group-rounded-linear', color 'bg-blue-50 text-blue-600' },
  ];
  return styles[index % styles.length];
};

function GenreRow({ genre, index }) {
  const [animes, setAnimes] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() = {
    const fetchAnimes = async () = {
      try {
         Stagger by 400ms per row to respect Jikan API rate limit (3 requestssecond)
        await new Promise(resolve = setTimeout(resolve, index  400));
        
        const res = await fetch(`httpsapi.jikan.moev4animegenres=${genre.mal_id}&limit=10&order_by=popularity`);
        let data = await res.json();
        
         Handle rate limiting fallback
        if (res.status === 429) {
           await new Promise(resolve = setTimeout(resolve, 2000));
           const retryRes = await fetch(`httpsapi.jikan.moev4animegenres=${genre.mal_id}&limit=10&order_by=popularity`);
           data = await retryRes.json();
        }
        
        setAnimes(data.data  []);
      } catch (error) {
        console.error(`Failed to fetch animes for ${genre.name}`, error);
      } finally {
        setLoading(false);
      }
    };
    fetchAnimes();
  }, [genre.mal_id, index]);

  const style = getGenreStyle(index);

  return (
    section className=flex flex-col
      div className=flex items-center justify-between mb-5
        div className=flex items-center gap-3
          div className={`flex items-center justify-center w-10 h-10 rounded-xl ${style.color} bg-opacity-50 border border-current10 shadow-sm`}
            iconify-icon icon={style.icon} width=20iconify-icon
          div
          h3 className=text-lg font-medium tracking-tight text-gray-900Explore {genre.name}h3
        div
        button className=text-sm font-medium text-gray-500 hovertext-gray-900 flex items-center gap-1 transition-colors group
          See All {genre.name}
          iconify-icon icon=solaralt-arrow-right-linear width=16 class=group-hovertranslate-x-0.5 transition-transformiconify-icon
        button
      div
      
      {loading  (
         div className=flex items-center justify-center h-32 text-sm text-gray-400 bg-gray-50 rounded-2xl border border-gray-100 border-dashed
           Loading {genre.name} anime...
         div
      )  (
        div className=grid grid-cols-5 gap-5
          {animes.map((anime) = (
            div key={anime.mal_id} className=flex flex-col group cursor-pointer
              div className=relative rounded-2xl overflow-hidden aspect-[43] mb-3 shadow-[0_4px_12px_rgb(0,0,0,0.03)] border border-gray-10050 bg-gray-100
                img 
                  src={anime.images.jpg.large_image_url  anime.images.jpg.image_url} 
                  alt={anime.title}
                  className=w-full h-full object-cover transition-transform duration-700 group-hoverscale-105
                
                { Play Button Overlay }
                div className=absolute bottom-3 left-3 w-8 h-8 rounded-full bg-white90 backdrop-blur shadow-sm flex items-center justify-center text-gray-800 opacity-90 group-hoveropacity-100 group-hoverscale-110 transition-all
                  iconify-icon icon=solarplay-bold width=14 class=ml-0.5iconify-icon
                div
              div
              h4 className=text-sm font-medium text-gray-900 tracking-tight truncate title={anime.title}{anime.title}h4
              p className=text-xs text-gray-500 mt-0.5{anime.year  'NA'} • {anime.episodes  `${anime.episodes} Episodes`  'Ongoing'}p
            div
          ))}
        div
      )}
    section
  );
}

function GenreSection() {
  const [genres, setGenres] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() = {
    const fetchGenres = async () = {
      try {
        const res = await fetch('httpsapi.jikan.moev4genresanime');
        const data = await res.json();
        
         Fetch top 10 genres to avoid overwhelming rate limits, while matching the user's minimum at 10 request.
        const topGenres = (data.data  []).filter(g = g.count  1000).slice(0, 10);
        setGenres(topGenres);
      } catch (error) {
        console.error(Failed to fetch genres, error);
      } finally {
        setLoading(false);
      }
    };
    fetchGenres();
  }, []);

  if (loading) {
    return (
      div className=mt-8 mb-8 flex items-center justify-center h-48 text-gray-500
        iconify-icon icon=solarspinner-linear width=24 class=animate-spin mr-2iconify-icon
        Loading categories...
      div
    );
  }

  return (
    div className=flex flex-col gap-12 mt-8 mb-8
      {genres.map((genre, idx) = (
        GenreRow key={genre.mal_id} genre={genre} index={idx} 
      ))}
    div
  );
}