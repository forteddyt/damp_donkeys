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

  genNewCode(event: any){
    // Make api call to code generator api here
    alert("New code is: ")
  }

  deleteCompany(event: any){
    alert("Deleted company (but not really yet)")
  }
}
