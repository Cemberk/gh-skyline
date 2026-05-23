// GitHub API module for fetching contribution data
// This module handles all GitHub API interactions

const GITHUB_API_URL = 'https://api.github.com/graphql';

/**
 * Fetches contribution data for a user and year from GitHub GraphQL API
 * @param {string} username - GitHub username
 * @param {number} year - Year to fetch contributions for
 * @param {string} token - Optional GitHub personal access token
 * @returns {Promise<Object>} Contribution data with weeks and days
 */
export async function fetchContributions(username, year, token = null) {
    const startDate = `${year}-01-01T00:00:00Z`;
    const endDate = `${year}-12-31T23:59:59Z`;

    const query = `
        query ContributionGraph($username: String!, $from: DateTime!, $to: DateTime!) {
            user(login: $username) {
                login
                contributionsCollection(from: $from, to: $to) {
                    contributionCalendar {
                        totalContributions
                        weeks {
                            contributionDays {
                                contributionCount
                                date
                            }
                        }
                    }
                }
            }
            rateLimit {
                limit
                remaining
                resetAt
            }
        }
    `;

    const variables = {
        username,
        from: startDate,
        to: endDate
    };

    const headers = {
        'Content-Type': 'application/json',
    };

    // Add authorization header if token is provided
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    try {
        const response = await fetch(GITHUB_API_URL, {
            method: 'POST',
            headers,
            body: JSON.stringify({ query, variables })
        });

        if (!response.ok) {
            throw new Error(`GitHub API error: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();

        if (data.errors) {
            throw new Error(`GraphQL error: ${data.errors.map(e => e.message).join(', ')}`);
        }

        if (!data.data.user) {
            throw new Error(`User "${username}" not found`);
        }

        return {
            contributions: data.data.user.contributionsCollection.contributionCalendar,
            rateLimit: data.data.rateLimit
        };
    } catch (error) {
        console.error('Error fetching contributions:', error);
        throw error;
    }
}

/**
 * Formats contribution data for WASM consumption
 * @param {Object} contributionCalendar - Contribution calendar from GitHub API
 * @returns {Array} Formatted contributions as 2D array [week][day]
 */
export function formatContributionsForWASM(contributionCalendar) {
    return contributionCalendar.weeks.map(week =>
        week.contributionDays.map(day => ({
            contributionCount: day.contributionCount,
            date: day.date
        }))
    );
}

/**
 * Formats rate limit information for display
 * @param {Object} rateLimit - Rate limit info from GitHub API
 * @returns {string} Formatted rate limit message
 */
export function formatRateLimit(rateLimit) {
    const resetDate = new Date(rateLimit.resetAt);
    const resetTime = resetDate.toLocaleTimeString();
    return `API Rate Limit: ${rateLimit.remaining}/${rateLimit.limit} remaining (resets at ${resetTime})`;
}
