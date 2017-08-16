class CreateSources < ActiveRecord::Migration[5.0]
	def change
		create_table :sources, id: false do |t|
			t.integer :source_id
			t.belongs_to :dataset, index: true
			t.text :source_name

			t.timestamps
		end
		execute "ALTER TABLE sources ADD PRIMARY KEY (source_id);"
	end
end
