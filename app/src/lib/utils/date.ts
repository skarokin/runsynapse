export function prettyPrintDate(date: string): string {
    // if over 1 week ago, return date. else, return 'x days ago' or 'x hours ago' or whatever
    const now = new Date();
    const thoughtDate = new Date(date);
    const diff = now.getTime() - thoughtDate.getTime();
    const oneWeek = 7 * 24 * 60 * 60 * 1000; // milliseconds in a week
    if (diff > oneWeek) {
        return thoughtDate.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
    }

    const options: Intl.DateTimeFormatOptions = {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: true,
        year: 'numeric',
    };
    return new Date(date).toLocaleString('en-US', options);
}