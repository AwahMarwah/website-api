package content_page

import content_page "website-api/model/content-page"

func (r *repo) SeedCmsPage() (err error) {
	return r.db.Create([]content_page.CmsPage{
		{
			Slug:  "about-us",
			Title: "About Us",
			Content: `
<p><strong>Website API</strong> adalah platform e-commerce dan CMS yang menyediakan solusi backend
yang aman, cepat, dan scalable.</p>
<p>Kami berfokus pada pengembangan sistem yang mudah dikembangkan dan handal.</p>
`,
			Status: true,
		},
		{
			Slug:  "terms-and-conditions",
			Title: "Terms & Conditions",
			Content: `
<p>Dengan menggunakan website ini, Anda menyetujui syarat dan ketentuan berikut:</p>
<ul>
  <li>Pengguna wajib memberikan data yang benar</li>
  <li>Dilarang menyalahgunakan sistem</li>
  <li>Kami berhak mengubah ketentuan sewaktu-waktu</li>
</ul>
`,
			Status: true,
		},
		{
			Slug:  "privacy-policy",
			Title: "Privacy Policy",
			Content: `
<p>Kami menghargai privasi pengguna dan berkomitmen melindungi data pribadi Anda.</p>
<p>Data yang dikumpulkan hanya digunakan untuk keperluan layanan.</p>
`,
			Status: true,
		},
	}).Error
}
