/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ContentOnDashboard {
  type?: string;
  count?: number;
}

export interface ActivityOnDate {
  date?: string;
  count?: number;
}

export interface Dashboard {
  numberOfNotes?: number;
  mostVisited?: ContentOnDashboard[];
  lastVisited?: ContentOnDashboard[];
  activityGraph?: ActivityOnDate[];
}

export interface Tag {
  typeOfParent?: string;
  value?: string;
  frequency?: number;
}

export interface PaginationMeta {
  range?: string;
  count?: number;
  hasNextPage?: boolean;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortDirection?: string;
}

export interface NoteInput {
  tags?: string[];
  title?: string;
  content?: string;
}

export interface NoteLite {
  shortId?: string;
  title?: string;
}

export interface NoteReferenceLink {
  source?: string;
  target?: string;
}

export interface NoteReferences {
  meta?: PaginationMeta;
  nodes?: NoteLite[];
  links?: NoteReferenceLink[];
}

export interface Note {
  shortId?: string;
  tags?: string[];
  title?: string;
  content?: string;
  createdAt?: string;
  createdBy?: User;
  updatedAt?: string;
  updatedBy?: User;
  hasDraft?: boolean;
  references?: NoteReferences;
}

export interface User {
  shortId?: string;
  fullname?: string;
  username?: string;
  email?: string;
  createdAt?: string;
  language?: string;
  theme?: string;
}

export interface UserActivity {
  clientId?: string;
  resourceType?: string;
  action?: string;
  objectId?: string;
  timestamp?: string;
}

export interface NoteSearchResult {
  meta?: PaginationMeta;
  items?: Note[];
}

export interface UserActivityResult {
  meta?: PaginationMeta;
  items?: UserActivity[];
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (securityData: SecurityDataType | null) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) => fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
    return keys
      .map((key) => (Array.isArray(query[key]) ? this.addArrayQueryParam(query, key) : this.addQueryParam(query, key)))
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string") ? JSON.stringify(input) : input,
    [ContentType.Text]: (input: any) => (input !== null && typeof input !== "string" ? JSON.stringify(input) : input),
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
            ? JSON.stringify(property)
            : `${property}`,
        );
        return formData;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(`${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`, {
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
      },
      signal: (cancelToken ? this.createAbortSignal(cancelToken) : requestParams.signal) || null,
      body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
    }).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title The Zula App
 * @version 1.0.0
 * @license MIT (https://opensource.org/license/mit/)
 * @contact <alperkilci@gmail.com>
 *
 * Organize your thoughts by writing notes, uploading files, saving bookmarks and searching easily on any device with a browser.
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  api = {
    /**
     * No description
     *
     * @name V1MeList
     * @summary Get the authenticated user
     * @request GET:/api/v1/me
     */
    v1MeList: (params: RequestParams = {}) =>
      this.request<User, any>({
        path: `/api/v1/me`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1UsersDetail
     * @summary Get a user by ID
     * @request GET:/api/v1/users/{shortId}
     */
    v1UsersDetail: (shortId: string, params: RequestParams = {}) =>
      this.request<User, any>({
        path: `/api/v1/users/${shortId}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1UsersActivityDetail
     * @summary Get user activity
     * @request GET:/api/v1/users/{shortId}/activity
     */
    v1UsersActivityDetail: (
      shortId: string,
      query?: {
        page?: number;
        pageSize?: number;
        sortBy?: string;
        sortDirection?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<UserActivityResult, any>({
        path: `/api/v1/users/${shortId}/activity`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1UsersInsightsDetail
     * @summary Get dashboard insights
     * @request GET:/api/v1/users/{shortId}/insights
     */
    v1UsersInsightsDetail: (shortId: string, params: RequestParams = {}) =>
      this.request<Dashboard, any>({
        path: `/api/v1/users/${shortId}/insights`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1TagsList
     * @summary Get tags
     * @request GET:/api/v1/tags
     */
    v1TagsList: (
      query?: {
        type?: string;
        q?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<Tag[], any>({
        path: `/api/v1/tags`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1NotesList
     * @summary List notes
     * @request GET:/api/v1/notes
     */
    v1NotesList: (
      query?: {
        q?: string;
        page?: number;
        pageSize?: number;
        sortBy?: string;
        sortDirection?: string;
        tags?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<NoteSearchResult, any>({
        path: `/api/v1/notes`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1NotesCreate
     * @summary Create a new note
     * @request POST:/api/v1/notes
     */
    v1NotesCreate: (data: NoteInput, params: RequestParams = {}) =>
      this.request<Note, any>({
        path: `/api/v1/notes`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1NotesDetail
     * @summary Get a note by shortId
     * @request GET:/api/v1/notes/{shortId}
     */
    v1NotesDetail: (
      shortId: string,
      query?: {
        loadDraft?: boolean;
        optOutTracking?: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<Note, any>({
        path: `/api/v1/notes/${shortId}`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1NotesUpdate
     * @summary Update a note by shortId
     * @request PUT:/api/v1/notes/{shortId}
     */
    v1NotesUpdate: (shortId: string, data: NoteInput, params: RequestParams = {}) =>
      this.request<boolean, any>({
        path: `/api/v1/notes/${shortId}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1NotesDelete
     * @summary Delete a note by shortId
     * @request DELETE:/api/v1/notes/{shortId}
     */
    v1NotesDelete: (shortId: string, params: RequestParams = {}) =>
      this.request<boolean, any>({
        path: `/api/v1/notes/${shortId}`,
        method: "DELETE",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1NoteDraftUpdate
     * @summary Save draft of a note by notes shortId
     * @request PUT:/api/v1/note/{shortId}/draft
     */
    v1NoteDraftUpdate: (shortId: string, data: NoteInput, params: RequestParams = {}) =>
      this.request<boolean, any>({
        path: `/api/v1/note/${shortId}/draft`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @name V1NoteDraftDelete
     * @summary Delete a notes draft by note shortId
     * @request DELETE:/api/v1/note/{shortId}/draft
     */
    v1NoteDraftDelete: (shortId: string, params: RequestParams = {}) =>
      this.request<boolean, any>({
        path: `/api/v1/note/${shortId}/draft`,
        method: "DELETE",
        format: "json",
        ...params,
      }),
  };
}
