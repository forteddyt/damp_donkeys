import { Component } from '@angular/core';
import { Router } from '@angular/router';


@Component({
  selector: 'app-select-company-tile',
  templateUrl: './select-company-tile.component.html',
  styleUrls: ['./select-company-tile.component.css']
})
export class SelectCompanyTileComponent  {
  companyName
  stateData
  constructor(private router: Router) { }

  reviewRedirect(event: any) {
  	// Append companyName to stateData
  	this.stateData.companyName = this.companyName
  	this.router.navigate(['/review'], { state: this.stateData })
  }
}
