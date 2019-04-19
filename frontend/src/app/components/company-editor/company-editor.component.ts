import { Component, AfterViewInit, ComponentFactoryResolver, ViewChild, ViewContainerRef } from '@angular/core';
import { Router} from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { CompanyNameTileComponent } from '../company-name-tile/company-name-tile.component';

@Component({
  selector: 'app-company-editor',
  templateUrl: './company-editor.component.html',
  styleUrls: ['./company-editor.component.css']
})
export class CompanyEditorComponent implements AfterViewInit {
	@ViewChild('companyInsert', { read: ViewContainerRef }) companyInsert: ViewContainerRef;
	stateData
	companyList

	constructor(private http: HttpClient, private router: Router, private vcr: ViewContainerRef, private cfr: ComponentFactoryResolver) {
		this.stateData = this.router.getCurrentNavigation().extras.state;
		if(this.stateData == null || this.stateData.jwt == null){
			//console.log("Invalid state has been passed (or not passed at all);"+this.stateData);
			this.router.navigateByUrl('/admin');
		}
	}

	async ngAfterViewInit() {
		const resp = await this.getCompanies();

		//console.log(resp.status)
		this.companyList = resp.body["company_list"];
		this.loadComponents()
	}

	getCompanies() {
		return this.http.get("https://csrcint.cs.vt.edu/api/company_list", {observe: 'response'}).toPromise();
	}

	loadComponents() {
		//console.log(this.companyList)
		const cFactory = this.cfr.resolveComponentFactory(CompanyNameTileComponent);
		this.companyInsert.clear();
		for(var i in this.companyList){
			const companyComponent = <CompanyNameTileComponent>this.companyInsert.createComponent(cFactory).instance;

			companyComponent.companyName = this.companyList[i];
			companyComponent.stateData = this.stateData;
		}
	}

	addRedirect(event: any) {
		// Append companyName to stateData
		this.router.navigateByUrl('/admin/companies/add', { state: this.stateData })
	}

	backRedirect(event: any){
		this.router.navigateByUrl('/admin/nav', { state: this.stateData })
	}
}
