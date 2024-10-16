export type PostSignInStartRequest = {
  redirect_to?: string;
};

export type PostSignInEmailStartRequest = PostSignInStartRequest & {
  email: string;
};

export type PostSignInEmailStartResponse = {
  sign_in_code_state: string;
  sign_in_code_id: string;
};

export type PostSignInOAuthStartRequest = PostSignInStartRequest;

export type PostSignInOAuthStartResponse = {
  auth_url: string;
};

export type PostSignInSAMLStartRequest = PostSignInStartRequest;

export type PostSignInSAMLStartResponse = Record<PropertyKey, never>;

export type PostSignInEndRequest = {
  state: string;
};

export type PostSignInEndResponse = {
  redirect_to?: string;
};

export type PostSignInEmailEndRequest = PostSignInEndRequest & {
  sign_in_code_id: string;
  sign_in_code_code: string;
};

export type PostSignInEmailEndResponse = PostSignInEndResponse;

export type PostSignInOAuthEndRequest = PostSignInEndRequest & {
  auth_result: string;
};

export type PostSignInOAuthEndResponse = PostSignInEndResponse;

export type PostSignInSAMLEndRequest = PostSignInEndRequest;

export type PostSignInSAMLEndResponse = PostSignInEndResponse;

export type PostSignOutRequest = Record<PropertyKey, never>;

export type PostSignOutResponse = Record<PropertyKey, never>;
