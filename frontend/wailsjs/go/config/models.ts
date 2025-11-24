export namespace config {
  export class ConfigItem {
    provider: string;
    name: string;
    fileName: string;
    path: string;
    scope: string;
    format: string;
    exists: boolean;

    static createFrom(source: any = {}) {
      return new ConfigItem(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.provider = source.provider;
      this.name = source.name;
      this.fileName = source.fileName;
      this.path = source.path;
      this.scope = source.scope;
      this.format = source.format;
      this.exists = source.exists;
    }
  }
}
