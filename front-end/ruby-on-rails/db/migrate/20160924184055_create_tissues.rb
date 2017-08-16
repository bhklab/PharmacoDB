class CreateTissues < ActiveRecord::Migration[5.0]
	def change
		create_table :tissues, id: false do |t|
			t.integer :tissue_id
			t.text :tissue_name

			t.timestamps
		end
		execute "ALTER TABLE tissues ADD PRIMARY KEY (tissue_id);"
	end
end
