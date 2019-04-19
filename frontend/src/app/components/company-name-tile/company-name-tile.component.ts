import { Component } from '@angular/core';
import { Router } from '@angular/router';


@Component({
  selector: 'app-company-name-tile',
  templateUrl: './company-name-tile.component.html',
  styleUrls: ['./company-name-tile.component.css']
})
export class CompanyNameTileComponent  {
  companyName
  stateData
  constructor(private router: Router) { }

  viewRedirect(event: any){
    this.stateData.companyName = this.companyName
    this.router.navigateByUrl('/admin/companies/view', { state: this.stateData })
  }
}
