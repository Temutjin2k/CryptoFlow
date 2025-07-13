const themeSwitch = document.getElementById('theme-switch');

const currentTheme = localStorage.getItem('theme') || 
                    (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');

if (currentTheme === 'dark') {
  document.documentElement.setAttribute('data-theme', 'dark');
  themeSwitch.checked = true;
}

themeSwitch.addEventListener('change', function(e) {
  if (e.target.checked) {
    document.documentElement.setAttribute('data-theme', 'dark');
    localStorage.setItem('theme', 'dark');
  } else {
    document.documentElement.setAttribute('data-theme', 'light');
    localStorage.setItem('theme', 'light');
  }
});

document.documentElement.classList.add('transition');

setTimeout(() => {
  document.documentElement.classList.remove('transition');
}, 300);

async function fetchLatestPrices() {
    const symbols = ["BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"];
    const exchanges = ["exchange1", "exchange2", "exchange3"];
    const container = document.getElementById("latest-prices");
    container.innerHTML = "";

    for (const exchange of exchanges) {
        for (const symbol of symbols) {
            try {
                const res = await fetch(`/prices/latest/${exchange}/${symbol}`);
                if (!res.ok) continue;
                const data = await res.json();
                const entry = document.createElement("div");
                entry.className = "price-entry";
                entry.innerHTML = `
                    <strong>${data.data.symbol}</strong> | ${data.data.exchange} — 
                    $${data.data.price.toFixed(2)} <br>
                    <small>${new Date(data.data.timestamp).toLocaleString()}</small>
                `;
                container.appendChild(entry);
            } catch (err) {
                console.error(`Failed for ${exchange}/${symbol}:`, err);
            }
        }
    }
}

async function fetchAggregated() {
    const metric = document.getElementById("metric-select").value;
    const symbol = document.getElementById("symbol-select").value;
    const exchange = document.getElementById("exchange-select").value;
    const period = document.getElementById("period-input").value;
    const container = document.getElementById("aggregated-result");

    let url;
    if (exchange && exchange !== "ALL") {
        url = `/prices/${metric}/${exchange}/${symbol}?period=${encodeURIComponent(period)}`;
    } else {
        url = `/prices/${metric}/${symbol}?period=${encodeURIComponent(period)}`;
    }

    try {
        const res = await fetch(url);
        const result = await res.json();

        let value;
        if (metric === "highest") value = result.data?.max;
        else if (metric === "lowest") value = result.data?.min;
        else if (metric === "average") value = result.data?.average;

        if (value === undefined) {
            container.innerHTML = `<p>No ${metric} data available for this period.</p>`;
            return;
        }

        const exchangeDisplay = result.data.exchange === "all" ? "ALL Exchanges" : result.data.exchange.toUpperCase();

        container.innerHTML = `
            <strong>Metric:</strong> ${metric.toUpperCase()}<br>
            <strong>Symbol:</strong> ${result.data.symbol}<br>
            <strong>Exchange:</strong> ${exchangeDisplay}<br>
            <strong>Period:</strong> ${result.period}<br>
            <strong>Value:</strong> ${value.toFixed(2)}<br>
            <strong>Timestamp:</strong> ${new Date(result.data.timestamp).toLocaleString()}
        `;
    } catch (e) {
        console.error("Fetch error:", e);
        container.innerHTML = "Failed to load aggregated data.";
    }
}

fetchLatestPrices();