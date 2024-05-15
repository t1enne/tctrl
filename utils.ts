export const API_BASE = "https://tpca-2d830a68c047.herokuapp.com/api";
export const AUTH_HEADERS = {
  Authorization: "",
};

export const req = async <T>(
  url: string,
  body: Record<string, unknown>,
): Promise<{ data: T }> => {
  const fullUrl = API_BASE + url;
  return await fetch(fullUrl, {
    method: "POST",
    credentials: "include",
    headers: {
      "User-Agent":
        "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
      "Accept": "application/json, text/plain, */*",
      "Accept-Language": "en-US,it-IT;q=0.8,it;q=0.5,en;q=0.3",
      "Content-Type": "application/json",
      "Sec-Fetch-Dest": "empty",
      "Sec-Fetch-Mode": "cors",
      "Sec-Fetch-Site": "cross-site",
      ...AUTH_HEADERS,
    },
    "referrer": "https://tpca.raintonic.com/",
    "mode": "cors",
    body: JSON.stringify(body),
  }).then((j) => j.json());
};
