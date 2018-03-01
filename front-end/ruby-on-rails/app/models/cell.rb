class Cell < ApplicationRecord
	has_many :source_cell_names
	has_many :cell_tissues
	has_many :experiments
	has_many :mol_cells
	has_many :dataset_cells
	has_many :datasets, through: :dataset_cells
	has_many :dose_responses, through: :experiments
	has_many :tissues, through: :cell_tissues
	has_many :sources, through: :source_cell_names
	has_one  :cellosauru
end