class CreateSourceTissueNames < ActiveRecord::Migration[5.0]
	def change
		create_table :source_tissue_names do |t|
			t.belongs_to :tissue, index: true
			t.belongs_to :source, index: true
			t.text :tissue_name

			t.timestamps
		end
	end
end
