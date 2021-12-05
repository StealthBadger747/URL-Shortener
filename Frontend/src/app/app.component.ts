import { Component } from '@angular/core';
import { ApiService } from '../services/api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'UrlShortener';
  enteredURL: string = "";
  shortenedURL: string = "";

  constructor(public apiService: ApiService) { }

  /**
   * Handles the click event for SHORTEN button.
   */
  public postShortenURL(): void {
    if (this.enteredURL.length === 0) {
      return;
    }

    this.apiService.createShortURL(this.enteredURL).then((value: any) => {
      this.shortenedURL = window.location.origin + "/" + value.short_url;
    });
  }
}
