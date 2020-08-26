export type Problem = {
    no: number
    title: string
    content: string
    testcases: any[]
    key?: number
    boilerplate: Boilerplate[]
}

export type Language =
    "javascript" | "java" | "golang" | "cpp"

export type Boilerplate = {
    lang: Language
    code: string
}

export type Result = {
    date: Date
    result: string
    passed: boolean
}
