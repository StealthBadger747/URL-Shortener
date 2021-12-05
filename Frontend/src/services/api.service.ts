import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  constructor(private http: HttpClient) { }

  /**
   * Makes the API call to get a shortened URL.
   * @param originalURL the original long url.
   * @returns a Promise of the request.
   */
  public createShortURL(originalURL: string): Promise<any> {
    const body = "";
    const params = new HttpParams()
      .set('url', originalURL);
    return new Promise(resolve => {
      this.http.post('/api/shorten_url', body, { params: params })
      .subscribe(response => {
        resolve(response);
      });
    });
  }
}
