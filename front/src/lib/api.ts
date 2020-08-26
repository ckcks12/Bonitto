import {Problem, Result} from "./type";

export async function APIGetProblems(): Promise<Problem[]> {
    try {
        return await fetch(`/api/problems`).then(a => a.json());
    } catch (e) {
        console.error(e);
        return [];
    }
}

export async function APIGetProblem(no: number): Promise<Problem> {
    try {
        return await fetch(`/api/problem/${no}`).then(a => a.json());
    } catch (e) {
        console.error(e);
        throw e;
    }
}

export async function APISubmit(no: number, id: string, code: string, lang: string) {
    try {
        await fetch(`/api/submit/${no}`, {
            body: JSON.stringify({id, code, lang}),
            headers: {
                "Content-Type": "application/json",
            },
            method: "POST",
        });
    } catch (e) {
        console.error(e);
        throw e;
    }
}

export async function APIGetResult(id: string, no: number): Promise<Result[]> {
    try {
        const rtn: Result[] = [];
        (await fetch(`/api/result/${id}/${no}`).then(a => a.json()))
            .forEach((results: Result[], scenarioIdx: number) => {
                results.forEach((result: Result, tcIdx: number) => {
                    const s = result.passed ? "✅" : "❌";
                    rtn.push({
                        date: new Date(results[0].date.toLocaleString()),
                        passed: true,
                        result: `Scenario #${scenarioIdx + 1} - Test #${tcIdx + 1} ${s} ${result.result}`,
                    });
                });
            });
        return rtn;
    } catch (e) {
        console.error(e);
        return [];
    }
}
