
export function filterEmptyValues(obj: Record<string, unknown>): Record<string, unknown> {
    return Object.entries(obj).reduce((acc, [key, value]) => {
        if (value !== null && value !== undefined && value !== "") {
            acc[key] = value;
        }
        return acc;
    }, {} as Record<string, unknown>);
}
