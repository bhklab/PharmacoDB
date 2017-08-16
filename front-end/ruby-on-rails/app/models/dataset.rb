class Dataset < ApplicationRecord
	has_many :sources
	has_many :experiments
	has_many :mol_cells
	has_many :dose_responses, through: :experiments
	has_many :targets, through: :experiments
end
