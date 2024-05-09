import React from "react";
import { HttpResponse } from "../types/Api";

type dataHookResponse<T, V> = [T, boolean, string | Error | null, (q: V) => Promise<void>];

export function useData<T, V>(fetchData: () => Promise<HttpResponse<T, unknown>>, query: V): dataHookResponse<T, V> {
    const [loading, setLoading] = React.useState(true);
    const [data, setData] = React.useState<T>();
    const [error, setError] = React.useState<string | null>(null);

    const fetchAndUpdate = async (query: V) => {
        try {
            setLoading(true);
            setError(null);



            let o = Object.keys(query)
            .filter((k) => query[k] != null)
            .reduce((a, k) => ({ ...a, [k]: query[k] }), {});

            const { data, status } = await fetchData(o)

            if (status === 200) {
                setData(data);
            } else {
                console.error(data);
                setError(data?.detail);
            }

        } catch (e: unknown) {
            console.error(e);

            setError(e as string);
        }
        setLoading(false);
    };

    return [data as T, loading, error, fetchAndUpdate];
}
