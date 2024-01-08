declare module '*.scss?inline' {
      const content: {[className: string]: string};
      export = content;
}